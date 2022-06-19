
# Splitwise system design 


Splitwise is an app for splitting expenses with your friends. It lets you and your friends add various bills and keep track of who owes who, and then it helps you to settle up with each other.

---

## Setup
- Run docker compose \
  ```docker compose up```
 
This will start 2 container:
- postgres: Port 5432
- webserver: 8080

---

## Endpoints
- Add user
  - Method: POST, Endpoint: 127.0.0.1:8080/user
  - request body: \
    `{
    "name":"vipul",
    "phone_number":"3139965135",
    "email_id":"s13s@example.com"
    }` \
  This will add new user in users table. `phone_number` and `email_id` must be unique.
  
- Add Group
  - Method: POST, Endpoint: 127.0.0.1:8080/group
  - request body \ `{
    "name":"testgroup",
    "userid":[1, 2, 3]
    }` \
  This will create new group. `userid` contains users belong to the group. Make sure user is present before
  creating new user.
- Create Expense
  - Method: POST, Endpoint: 127.0.0.1:8080/expense
  - request body \ `{
    "paid_by": 3,
    "category":"shop",
    "Amount":1400,
    "paid_to":[1]
    }}`.
  This will create new expense in the database. `paid_by` is the user who paid for the expense and `paid_to`
  are the user who are owned to `paid_by`. `paid_to` and `paid_by` are userids.
  
- Add Group Expense
    - Method: POST, Endpoint: 127.0.0.1:8080/group_expense
    - request body \ `{
      "paid_by":2,
      "category":"books",
      "Amount":2400,
      "gid":1
      }`
    This is same as above call instead of mentioning paid_to user we are sending the group using `gid`. And the
    expense will be divided among group equally.
- Modify Expense 
    - Method: PATCH, Endpoint: 127.0.0.1:8080/expense
    - request body \ `{
      "id":4,
      "amount":1600
      }`
    This will modify the amount of existing expense. Here `id` is expense id that we got when we created new expense.
- View Expense
    - Method: GET, Endpoint: 127.0.0.1:8080/expense 
    - request body \ `{"id":1}`
    This will show the expense corresponding to given `id`.
- Settle Expense
  - Method: POST, Endpoint: 127.0.0.1:8080/settle  
  - request body \ `{"id":1}`
    This will settle the expense corresponding to given `id`.
- Summary
  - Method: Get, Endpoint: 127.0.0.1:8080/summary
  - request body \ `{"id":1}`
    This will fetch the balance corresponding to given `id`. Here `id` is user id that we got when we added the user.
- All Unsettled
  - Method: Get, Endpoint: 127.0.0.1:8080/all_unsettled
  - This will fetch all unsettled expense.



