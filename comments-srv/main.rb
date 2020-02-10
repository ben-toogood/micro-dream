require 'micro'
require_relative 'sdk/ruby/comments_service.rb'

class Service < Micro::Service
    def create_comment(req, rsp)
        # ctx.Metadata.Get("FooBar")
    end
end

CommentsService.RegisterCommentsServiceHandler(Service.new)