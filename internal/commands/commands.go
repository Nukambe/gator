package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	cfg "github.com/Nukambe/gator/internal/config"
	"github.com/Nukambe/gator/internal/database"
	publish2 "github.com/Nukambe/gator/internal/publish"
	"github.com/Nukambe/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type State struct {
	Db     *database.Queries
	Config *cfg.Config
}

type Command struct {
	Name string
	Args []string
}

type commandHandler func(*State, Command) error

type Commands struct {
	cmds map[string]commandHandler
}

func (c *Commands) register(name string, handler commandHandler) {
	c.cmds[name] = handler
}

func (c *Commands) Run(s *State, cmd Command) error {
	if s == nil {
		return fmt.Errorf("state is nil")
	}
	handler, ok := c.cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("command '%s' does not exist", cmd.Name)
	}
	if err := handler(s, cmd); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}
	return nil
}

func (s *State) getCurrentUser() (database.User, error) {
	user, errUser := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if errUser != nil {
		return database.User{}, fmt.Errorf("unable to retrieve user: %w", errUser)
	}
	return user, nil
}

type checkLoginSignature func(s *State, cmd Command, user database.User) error

func middlewareLoggedIn(handler checkLoginSignature) commandHandler {
	return func(s *State, cmd Command) error {
		user, errUser := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if errUser != nil {
			return fmt.Errorf("unable to retrieve user: %w", errUser)
		}
		return handler(s, cmd, user)
	}
}

func InitCommands() Commands {
	commands := Commands{cmds: map[string]commandHandler{}}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handleUsers)
	commands.register("agg", handleAgg)
	commands.register("addfeed", middlewareLoggedIn(handleAddFeed))
	commands.register("feeds", handleFeeds)
	commands.register("follow", middlewareLoggedIn(handleFollow))
	commands.register("following", middlewareLoggedIn(handleFollowing))
	commands.register("unfollow", middlewareLoggedIn(handleUnfollow))
	commands.register("browse", middlewareLoggedIn(handleBrowse))
	return commands
}

// login
func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login expects a username")
	}
	// check if user exists
	user, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to login: %w", err)
	}
	// login
	if err = s.Config.SetUser(user.Name); err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	// show message if the login command was used
	if cmd.Name == "login" {
		fmt.Println("logged in as:", user.Name)
	}
	return nil
}

// register
func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register expects a name")
	}

	// create user
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	// login with new user
	if err = handlerLogin(s, cmd); err != nil {
		return fmt.Errorf("unable to login after register: %w", err)
	}
	fmt.Printf("User created: %v\n", user)
	return nil
}

// reset
func handlerReset(s *State, cmd Command) error {
	if err := s.Db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("unable to delete users: %w", err)
	}
	fmt.Println("deleted all users")
	return nil
}

// list users
func handleUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}
	fmt.Println("registered users:")
	for _, user := range users {
		fmt.Printf("	* %s", user.Name)
		if s.Config.CurrentUserName == user.Name {
			fmt.Println(" (current)")
		} else {
			fmt.Println()
		}
	}
	return nil
}

// aggregate
func handleAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("agg expects one arg")
	}
	duration := cmd.Args[0]

	tick, errTime := time.ParseDuration(duration)
	if errTime != nil {
		return fmt.Errorf("unable to parse time: %w", errTime)
	}

	fmt.Printf("Collecting feeds every %v...\n", tick)

	ticker := time.NewTicker(tick)
	for ; ; <-ticker.C {
		err := s.scrapeFeeds()
		if err != nil {
			return fmt.Errorf("unable to scrape feed: %w", err)
		}
	}
}

func (s *State) scrapeFeeds() error {
	// get next feed
	nextFeed, errNext := s.Db.GetNextFeedToFetch(context.Background())
	if errNext != nil {
		return fmt.Errorf("unable to get next feed: %w", errNext)
	}

	// mark feed
	errMark := s.Db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if errMark != nil {
		return fmt.Errorf("unable to make feed fetched: %w", errMark)
	}

	// fetch feed
	feed, errFetch := rss.FetchFeed(context.Background(), nextFeed.Url)
	if errFetch != nil {
		return fmt.Errorf("unable to fetch feed: %w", errFetch)
	}

	// save posts
	for _, post := range (*feed).Channel.Item {
		// verify nullable description
		description := sql.NullString{}
		description.String = post.Description
		description.Valid = post.Description != ""

		// verify nullable publish date
		publish := publish2.ParsePubDate(post.PubDate)

		// save post
		if _, err := s.Db.CreatePost(context.Background(), database.CreatePostParams{
			Title:       post.Title,
			Url:         post.Link,
			Description: description,
			PublishedAt: publish,
			FeedID:      nextFeed.ID,
		}); err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					continue
				}
			}
			return fmt.Errorf("unable to create post: %w", err)
		}
	}
	return nil
}

// add feed
func handleAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("'add feed' requires two arguments")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	// get feed
	_, errFeed := rss.FetchFeed(context.Background(), url)
	if errFeed != nil {
		return fmt.Errorf("unable to fetch feed: %w", errFeed)
	}

	// save feed
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed: %w", err)
	}

	// follow feed
	if err = handleFollow(s, Command{
		Name: cmd.Name,
		Args: []string{feed.Url},
	}, user); err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

// list feeds
func handleFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to retrieve feeds: %w", err)
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		fmt.Printf("	* '%s' %s - %s\n", feed.FeedName, feed.FeedUrl, feed.UserName)
	}
	return nil
}

// follow feed
func handleFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow requires one argument")
	}
	url := cmd.Args[0]

	// get feed
	feed, errFeed := s.Db.GetFeedByUrl(context.Background(), url)
	if errFeed != nil {
		return fmt.Errorf("unable to retrieve feed by url: %w", errFeed)
	}

	// save follow
	follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %w", err)
	}

	fmt.Printf("%s followed %s\n", follow.UserName, follow.FeedName)
	return nil
}

// list following
func handleFollowing(s *State, cmd Command, user database.User) error {
	// get follows
	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("unable to retreive feed_follows: %w", err)
	}

	fmt.Println("Following:")
	for _, follow := range follows {
		fmt.Printf("	* %s\n", follow.FeedName)
	}
	return nil
}

// unfollow
func handleUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow requires one argument")
	}
	url := cmd.Args[0]

	err := s.Db.DeleteFeedFollowByUserIdAndURL(context.Background(), database.DeleteFeedFollowByUserIdAndURLParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return fmt.Errorf("unable to delete follow record: %w", err)
	}
	return nil
}

// browse
func handleBrowse(s *State, cmd Command, user database.User) error {
	limit := int32(2)
	if len(cmd.Args) > 0 {
		if num64, err := strconv.ParseInt(cmd.Args[0], 10, 32); err == nil {
			limit = int32(num64)
		} else {
			return fmt.Errorf("unable to parse '%s', using 2", cmd.Args[0])
		}
	}

	// get posts
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("unable to retrieve posts: %w", err)
	}

	// print posts
	for _, post := range posts {
		fmt.Println(post)
	}
	return nil
}
