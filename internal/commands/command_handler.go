package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/TheHawk24/gator/internal/config"
	"github.com/TheHawk24/gator/internal/database"
	"github.com/TheHawk24/gator/internal/rss"
	"github.com/google/uuid"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands_all map[string]func(*State, Command) error
}

// Run a command handler
func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Commands_all[cmd.Name]
	if !ok {
		fmt.Println("Command not found")
		return errors.New("Command not found")
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands_all[name] = f
}

func HandlerLogin(s *State, cmd Command) error {

	if len(cmd.Args) == 0 {
		err := fmt.Sprintf("%v expects a single argument, a username", cmd.Name)
		fmt.Println(err)
		return errors.New(err)
	}

	name := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("User does not exist\n")
		return err
	}

	s.Config.Current_Username = cmd.Args[0]
	s.Config.SetUser()
	fmt.Printf("User %v is logged in\n", user.Name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) == 0 {
		err := fmt.Sprintf("%v expects a single argument, a username", cmd.Name)
		fmt.Println(err)
		return errors.New(err)
	}

	var register_user database.CreateUserParams
	register_user.ID = uuid.New()
	register_user.CreatedAt = time.Now()
	register_user.UpdatedAt = time.Now()
	register_user.Name = cmd.Args[0]

	user, err := s.Db.CreateUser(context.Background(), register_user)
	if err != nil {
		fmt.Printf("User already exists\n")
		return err
	}

	s.Config.Current_Username = cmd.Args[0]
	s.Config.SetUser()

	fmt.Printf("User %v successfully created\n", cmd.Args[0])
	fmt.Println("ID: ", user.ID)
	fmt.Println("Name: ", user.Name)
	fmt.Println("Created At: ", user.CreatedAt)
	fmt.Println("Updated At: ", user.UpdatedAt)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("Failed to delete users")
		return err
	}

	fmt.Println("Successfully deleted all users")
	return nil
}

func HandlerListUsers(s *State, cmd Command) error {

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		log.Println("Failed to retrieve all users")
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users found")
	}

	name := s.Config.Current_Username

	for _, user := range users {
		if user.Name == name {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Println("*", user.Name)
		}
	}

	return nil
}

func display_feed(feed *rss.RSSFeed) {

	fmt.Println("Title:", feed.Channel.Title)
	fmt.Println("Link:", feed.Channel.Link)
	fmt.Println("Description:", feed.Channel.Description)

	for _, v := range feed.Channel.Item {
		fmt.Println("-------------------------------------------------------------------------------------------")
		fmt.Println(" - Title:", v.Title)
		fmt.Println(" - Link:", v.Link)
		fmt.Println(" - PubDate:", v.PubDate)
		fmt.Println(" - Description:", v.Description)
		fmt.Println("-------------------------------------------------------------------------------------------")
	}

}

func HandlerAgg(s *State, cmd Command) error {

	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		log.Println(err)
		return err
	}

	display_feed(feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {

	if len(cmd.Args) < 2 {
		err := fmt.Sprintf("%v expects a two arguments, a feed name and url", cmd.Name)
		fmt.Println(err)
		return errors.New(err)
	}

	//Check if user exists
	current_user := s.Config.Current_Username
	user_info, err := s.Db.GetUser(context.Background(), current_user)
	if err != nil {
		log.Println("Cannot add feed for uknown user")
		return err
	}

	feed_name := cmd.Args[0]
	url := cmd.Args[1]

	//Add feed for user
	var database_feed database.CreateFeedParams
	database_feed.ID = uuid.New()
	database_feed.Name = feed_name
	database_feed.CreatedAt = time.Now()
	database_feed.UpdatedAt = time.Now()
	database_feed.Url = url
	database_feed.UserID = user_info.ID

	feed, err := s.Db.CreateFeed(context.Background(), database_feed)
	if err != nil {
		log.Println("Failed to add feed")
		return err
	}

	//FOLLOW FEED

	var feedFollow database.CreateFeedFollowsParams
	feedFollow.ID = uuid.New()
	feedFollow.CreatedAt = time.Now()
	feedFollow.UpdatedAt = time.Now()
	feedFollow.UserID = user_info.ID
	feedFollow.FeedID = feed.ID
	_, err = s.Db.CreateFeedFollows(context.Background(), feedFollow)
	if err != nil {
		log.Println("Failed to follow feed")
	}

	return nil

}

func HandlerFeeds(s *State, cmd Command) error {

	//Fetch feeds from database
	feeds, err := s.Db.GetUsersFeeds(context.Background())
	if err != nil {
		log.Println("Failed to fetch all feeds")
		return err
	}

	if len(feeds) == 0 {
		err := fmt.Sprintf("No feeds founds")
		log.Println(err)
		return fmt.Errorf(err)
	}

	for _, v := range feeds {
		name := v.Name
		url := v.Url
		username := v.Name_2
		fmt.Println("--------------------------------------------------------")
		fmt.Printf("Feed Name: %v\n", name)
		fmt.Printf("URL: %v\n", url)
		fmt.Printf("User name: %v\n", username)
		fmt.Println("--------------------------------------------------------")
	}

	return nil
}

func HandlerFollow(s *State, cmd Command) error {

	if len(cmd.Args) == 0 {
		err := fmt.Errorf("Command follow expects a single argument, a url")
		log.Println(err)
		return err
	}

	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), url)
	if err != nil {
		log.Println("This feed does not exist")
		return err
	}

	username := s.Config.Current_Username
	user_info, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		log.Println("User does not exist")
		return err
	}

	//Get user id and feed Id
	userid := user_info.ID
	feedID := feed.ID

	var feedFollow database.CreateFeedFollowsParams
	feedFollow.ID = uuid.New()
	feedFollow.CreatedAt = time.Now()
	feedFollow.UpdatedAt = time.Now()
	feedFollow.UserID = userid
	feedFollow.FeedID = feedID
	user_feed_follow, err := s.Db.CreateFeedFollows(context.Background(), feedFollow)
	if err != nil {
		log.Println("Failed to follow feed")
	}

	fmt.Printf("UserName: %v\n", user_feed_follow.UserName)
	fmt.Printf("FeedName: %v\n", user_feed_follow.FeedName)

	return nil
}

func HandlerFollowing(s *State, cmd Command) error {

	username := s.Config.Current_Username
	user_info, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		log.Println("User does not exist")
		return err
	}

	allfeeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user_info.ID)
	if err != nil {
		log.Println("Failed to get feed follows for user")
		return err
	}

	fmt.Printf("%v follows this feeds\n", user_info.Name)
	for _, v := range allfeeds {
		fmt.Println(v.FeedName)
	}

	return nil
}
