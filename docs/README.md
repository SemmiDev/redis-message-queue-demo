- Endpoints
  - create post
  - get post

- Admin
  - When admin create a post, will enqueue to `email` queue

- Subscriber users
  - When user click on link in email, will enqueue to `analytic` queue
  - When user read the post, will enqueue to `analytic` queue

- Email Worker
  - Notify to all subscribers with sending email with Urchin Tracking Module (UTM) parameters (for analytics) when admin create a post

- Analytic Worker
  - Count user read post
  - Count user read post by UTM parameters