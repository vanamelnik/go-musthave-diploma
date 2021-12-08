# gophermart

Online shop loyalty system.

The system is an HTTP API with the following business logic requirements:
- registration, authentication and authorization of users;
- receiving order numbers from registered users;
- accounting and maintaining a list of transferred order numbers of a registered user;
- registration and maintenance of the bonus account of the registered user;
- verification of accepted order numbers through the loyalty points calculation system;
- accrual for each suitable order number of the required reward to the user's loyalty account.

### Paths:
* `POST /api/user/register` - user registration;
* `POST /api/user/login` - user authentication;
* `POST /api/user/orders` - loading the order number by the user for calculation;
* `GET /api/user/orders` - getting a list of order numbers uploaded by the user, their processing statuses and information about charges;
* `GET /api/user/balance` - getting the current account balance of the user's bonus points;
* `POST /api/user/balance/withdraw` - a request to withdraw points from a bonus account to pay for a new order;
* `GET /api/user/balance/withdrawals` - receiving information about the withdrawal of funds from the bonus account by the user.