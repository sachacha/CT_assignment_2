Different things to know about this program.
I'm using firebase as NoSQL solution.
I'm using the auto generated IDs by firebase as IDs for my webhooks.
In case you don't know (just in case...), use Postman to test the different endpoints, even the GET ones 
(some browser doesn't make the difference between .../webhooks and .../webhooks/).
I'm sending the parameters used to the webhooks, not the content of those parameters, for security issues.
