package main

import (
	"context"
	"net/http"

	foobar "github.com/micro/dream/foobar-srv/sdk/go"
	sdk "github.com/micro/dream/posts-srv/sdk/go"
	"github.com/micro/micro.go/store"
)

func main() {
	handler := new(Handler)

	// If the developer needs to initialize any custom packages
	// they can do so here using os.Getenv. We should not be
	// providing any config tools for the developer as this is
	// outside our domain and there are existing solutions in place.

	// There should be no initialization required for micro - it
	// just works. The only exception is registering for messages
	// from the micro broker. To subscribe to a message, you
	// leverage the SDK of the service which publishes the message,
	// ensuring the messages are decoded to the correct type, and
	// then optionally provide a type. The default is that if srv
	// foobar publishes a message, only one instance of the
	// posts-srv will consume the message, however the developer
	// may want to override this so all instances can consume the
	// message.
	foobar.RegisterNewUserSubscriber(handler.handleFooBar, broker.TypeOneOfMany)

	// RegisterPostsServiceHandler will initialize a service, run
	// it, registerto service discovery and then wait for a cancel
	// command. If there is an error, the error will be logged and
	// the program will exit with a non-zero code.
	//
	// RegisterPostsServiceHandler takes a struct which impements
	// the sdk and no arguments for options. As we're building
	// platform first, the service shouldn't set any dependancies
	// in code. This could be done in the .microservice file
	sdk.RegisterPostsServiceHandler(handler)
}

// Handler implements the service definition
type Handler struct{}

// handleFooBar is an example of how one service would consume
// messages published via another service via the broker. The
// events are statically typed in the sdk. It may be that we
// can't find a way to do this using proto files and need to use
// another method of service definition
func (h *Handler) handleFooBar(foo foobar.Foo) error {
	// if an error is returned then the message is not marked as
	// recieved and it's then passed to another instance or
	// retried.

	// example of how we can call another service. Note, we don't
	// initalize a client and then perform the request, but the
	// 'GetFoo' method is exported by the package itself. This
	// goes back to the principle of
	resp, err := foobar.GetFoo()

	return nil
}

// Create a new post
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// Write to DB. The ORM needs consideration however the store
	// should not need to have a pointer reference in handler, we
	// can use a var default and have the package provide accessor
	// methods. This default can be overriden when it is deployed
	// using a auto-generated plugins.go file which the developer
	// will never see.
	//
	// Ideally we'd have an ORM which allows for validations and
	// leverages micro errors to propagate these effectively in
	// a way which useful to a frontend client, e.g. a key value
	// map which can be used to highlight invalid input to a web
	// user.
	user, err := store.Table("foo").Create(req.User)
	if err != nil {
		return err
	}

	// Publish message to the broker. Async messaging is a first
	// class citizen in micro, and the message types are clearly
	// defined in the SDK
	sdk.PublishNewUser(&sdk.User{
		Foo: user.Foo,
	})

	// Setting the response status code would be amazing since
	// this is commonly used by frontend apps. Not all response
	// codes mean an error. Possibly leverage http pkg to allow
	// users to say StatusCreated instead of 201.
	rsp.StatusCode = http.StatusCreated

	return nil
}
