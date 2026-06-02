# Finance Calculator

**A REST API with useful Savings and Loans calculations.**


## Technologies Used

- Go
- PostgreSQL

## Getting Started

### Prerequisites

Before running the project, make sure you have the following installed:

- [Go](https://go.dev/doc/install) 
- [PostgreSQL](https://www.postgresql.org/download/)
- [Goose](https://pkg.go.dev/github.com/pressly/goose/v3#section-readme) (Database migration tool)

---

### Installation

Clone the repository:

```bash
git clone git@github.com:Mr-Rafael/finance-calculator.git
cd finance-calculator
```

---

### Configure the Database

Create a PostgreSQL database:

```sql
CREATE DATABASE your_database_name;
```

Set the database connection environment variables as needed by the project.

Example:

```bash
ALLOWED_ORIGIN=http://localhost:5173
POSTGRES_CONNECTION_STRING=postgres://<username>:<password>@localhost:5432/finance_calculator?sslmode=disable
ACCESS_SECRET=DEVENVIRONMENTSECRET
REFRESH_SECRET=DEVENVREFRESHSECRET
ENV=develop
```

- **ALLOWED_ORIGIN**: Used for CORS. Necessary when you're running both the server and a client on the same computer.
- **POSTGRES_CONNECTION_STRING**: The user, password and address of the PostgreSQL database you're running.
- **ACCESS_SECRET**: The secret that will be used to sign Access Tokens.
- **REFRESH_SECRET**: The secret that will be used to sign Refresh Tokens.
- **ENV**: If set to "production", the Refresh Token cookie will be set to Secure, and only be sent via HTTPS.

---

### Installing Goose

If Goose is not installed:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

### Run Database Migrations

Apply the migrations using Goose:

```bash
goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=your_database_name sslmode=disable" up
```

---

### Install Dependencies

```bash
go mod tidy
```

---

### Run the Project

```bash
go run ./cmd/server/main.go
```

---

### Notes

- Make sure PostgreSQL is running before starting the application.
- Ensure the database credentials match your local setup.

## API Endpoints

| route               | description                                          
|----------------------|-----------------------------------------------------
| <kbd>GET /api/healthz</kbd>     | Check if server is running.
| <kbd>POST /app/users/create</kbd>     | Create a new user.
| <kbd>POST /app/login</kbd>     | Login user.
| <kbd>POST /app/refresh</kbd>     | Refresh the access token.
| <kbd>POST /app/savings/calculate</kbd>     | Generate a Savings Plan without saving.
| <kbd>POST /app/loans/calculate</kbd>     | Calculate a Loan Payment Plan without saving.
| <kbd>POST /app/savings/save</kbd>     | Calculate and save a Savings Plan.
| <kbd>POST /app/loans/save</kbd>     | Calculate and save a Loan Payment Plan.
| <kbd>GET /app/savings/list</kbd>     | List a User's saved Savings Plans. 
| <kbd>GET /app/loans/list</kbd>     | List a User's saved Loan Payment Plans.
| <kbd>GET /app/savings/{id}</kbd>     | Get a previously saved Savings Plan.
| <kbd>GET /app/loans/{id}</kbd>     | Get a previously saved Loan Payment Plan.
| <kbd>PATCH /app/savings/{id}</kbd>     | Update and recalculate a Savings Plan.
| <kbd>PATCH /app/loans/{id}</kbd>     | Update and recalculate a Loan Payment Plan.
| <kbd>DELETE /app/savings/{id}</kbd>     | Delete a Savings Plan.
| <kbd>DELETE /app/loans/{id}</kbd>     | Delete a Loan Payment Plan.

## `GET /api/healthz`

Short description of what the endpoint does.

---

### Request

### URL

```http
METHOD /endpoint
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | No | Bearer token if authentication is required |
| Content-Type | Yes | Usually `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | Resource identifier |

### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| limit | integer | No | Number of items to return |

### Request Body

```json
{
  "field1": "value",
  "field2": 123
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| field1 | string | Yes | Example string field |
| field2 | integer | No | Example numeric field |

---


### Success Response

**Status Code:** `200 OK`

```json
{
  "id": "123",
  "field1": "value",
  "created_at": "2026-05-28T12:00:00Z"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | Resource ID |
| field1 | string | Example field |
| created_at | string | ISO 8601 timestamp |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "invalid request body"
}
```

### `404 Not Found`

```json
{
  "error": "resource not found"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```

---


| route               | description                                          
|----------------------|-----------------------------------------------------
| <kbd>GET /api/healthz</kbd>     | Check if server is running.
| <kbd>POST /app/users/create</kbd>     | Create a new user.
| <kbd>POST /app/login</kbd>     | Login user.
| <kbd>POST /app/refresh</kbd>     | Refresh the access token.
| <kbd>POST /app/savings/calculate</kbd>     | Generate a Savings Plan without saving.
| <kbd>POST /app/loans/calculate</kbd>     | Calculate a Loan Payment Plan without saving.
| <kbd>POST /app/savings/save</kbd>     | Calculate and save a Savings Plan.
| <kbd>POST /app/loans/save</kbd>     | Calculate and save a Loan Payment Plan.
| <kbd>GET /app/savings/list</kbd>     | List a User's saved Savings Plans. 
| <kbd>GET /app/loans/list</kbd>     | List a User's saved Loan Payment Plans.
| <kbd>GET /app/savings/{id}</kbd>     | Get a previously saved Savings Plan.
| <kbd>GET /app/loans/{id}</kbd>     | Get a previously saved Loan Payment Plan.
| <kbd>PATCH /app/savings/{id}</kbd>     | Update and recalculate a Savings Plan.
| <kbd>PATCH /app/loans/{id}</kbd>     | Update and recalculate a Loan Payment Plan.
| <kbd>DELETE /app/savings/{id}</kbd>     | Delete a Savings Plan.
| <kbd>DELETE /app/loans/{id}</kbd>     | Delete a Loan Payment Plan.
## `GET /api/healthz`

Check if the server is running.

### URL

```http
GET /api/healthz
```
---


### Success Response

**Status Code:** `200 OK`

```text
OK
```
---
## `POST /app/users/create`

Create a new user.

---

### URL

```http
POST /app/users/create
```

### Headers

| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | `application/json` |

### Request Body

```json
{
	"email":"user@mail.com",
    "password":"password",
    "username":"User Name"
}
```

### Request Fields


| Parameter | Type | Required | Description |
|---|---|---|---|
| email | string | Yes | Email of the user to create. Must be unique to the user. |
| password | string | Yes | Password for the user. |
| username | string | Yes | Display name for the user. |

---

### Response

### Success Response

**Status Code:** `201 Created`

```json
{
    "ID": "fa3ad421-c507-4d3c-bd77-e7918a67aaae",
    "Email": "user@mail.com",
    "Username": "User Name",
    "CreatedAt": "1970-01-01T21:39:22.159435"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| ID | string | UUID of the created user. |
| Email | string | User's email. |
| Username | string | User's Display Name. |
| CreatedAt | string | ISO 8601 timestamp |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "invalid request body"
}
```
### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST app/login`

Login to receive an access and refresh token. Access tokens are necessary to access some endpoints.

---

### URL

```http
POST /app/login
```

### Headers

| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | `application/json` |

### Request Body

```json
{
  "email":"user@mail.com",
  "password":"password"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| email | string | Yes | User's email. |
| password | string | Yes | User's password. |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "70fe1047-9c22-42e8-baac-64772c5c475a",
    "email": "user@mail.com",
    "username": "Test",
    "access_token": "<access token here>"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | User's UUID (version 4) |
| email | string | User's email |
| username | string | User's display name |
| access_token | string | JWT token to access authenticated endpoints |

---

### Error Responses

### `401 Unauthorized`

```json
{
  "error": "invalid request body"
}
```
---
## `POST /app/refresh`

Obtain a new access token, if the old one has expired.

---

### URL

```http
POST /app/refresh
```

### Headers

| Header | Required | Description |
|---|---|---|
| Cookie | Yes | The cookie should contain the refresh token obtained from the Login endpoint. |

### Request Body

No body. All required information is on the cookie.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "access_token": "[access token here]"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| access_token | string | A new, valid access token to access authenticated endpoints. |

---

### Error Responses

### `401 Unauthorized`

```json
{
    "error": "missing refresh token"
}
```
## `POST /app/savings/calculate`

Calculate a Savings Plan without saving it.

---

### URL

```http
POST /app/savings/calculate
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
	"startingCapital": 700000,
	"yearlyInterestRate": "4.75",
    "interestRateType": "APR",
	"monthlyContribution": 15000,
	"durationYears": 1,
    "taxRate": "5",
    "yearlyInflationRate": "6",
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| startingCapital | integer | Yes | How much money was deposited at the start of the term, in cents. For example, startingCapital = 100 would mean $1. |
| yearlyInterestRate | string | Yes | The yearly interest rate for the savings plan. Send as a percent. For example, "6.25" would be a 6.25% interest rate. |
| startDate | integer | Yes | The start date of the savings plan. |
| durationYears | integer | Yes | The term you want to calculate in years. 1 means "calculate the savings plan for 1 year".  ISO 8601|
| interestRateType | string | No | Send "APR" or "APY", depending on the type of interest rate. If empty, it defaults to APY. |
| monthlyContribution | integer | No | The monthly deposits that will be made (if any). Defaults to 0 if not in the request. The amount is in cents (e.g. 15000 means $150) |
| taxRate | string | No | The tax rate paid on returns. Send as a percent. For example, "5" means a 5% tax rate. Defaults to 0% if not in the request. |
| yearlyInflationRate | string | No | The yearly inflation rate, used for rate of return calculations. Send as a percent. For example, "6" means a 6% yearly inflation rate. Defaults to 0% if not in the request..  |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "monthlyInterestRate": "0.4074123784",
    "totalEarnings": 500000,
    "totalDeposited": 10000000,
    "rateOfReturn": "5",
    "inflationAdjustedROR": "5",
    "plan": [
        {
            "date": "2026-03-01T00:00:00Z",
            "interest": 40741,
            "tax": 0,
            "contribution": 0,
            "increase": 40741,
            "capital": 10040741
        },
        {
            "date": "2026-03-01T00:00:00Z",
            "interest": 40741,
            "tax": 0,
            "contribution": 100,
            "increase": 40741,
            "capital": 10040741
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| monthlyInterestRate | string | The monthly interest rate that was used on the calculations. |
| totalEarnings | integer | The total amount in cents that you earned in interest. |
| totalDeposited | integer | The total amount deposited in cents. Includes the initial deposit and monthly deposits. |
| rateOfReturn | string | The total in the account at the end of the term, divided by the total deposits made. A measure of how much return was made. The value is a percent (e.g. rateOfReturn = "5" means a 5% rate of return). |
| inflationAdjustedROR | string | The rate of return divided by the total inflation over the term. The value is a percent (e.g. rateOfReturn = "5" means a 5% rate of return). |
| plan | array of plan statuses | An array of monthly statuses of the Savings Plan. |

**Plan Status Fields**
| Field | Type | Description |
|---|---|---|
| date | string | The date of this monthly status. ISO 8601. |
| interest | integer | The interest earned this month in cents. |
| tax | integer | The tax paid this month in cents. |
| contribution | integer | The deposit made this month in cents. |
| increase | integer | The increase in savings at the end of this month, in cents. Includes deposits and interest earnings minus taxes. |
| capital | integer | The total money in the account at the end of this month, in cents. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message dependimg on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST /app/loans/calculate`

Calculate a Loan Payment Plan without saving it.

---

### URL

```http
POST /app/loans/calculate
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
	"startingPrincipal": 10000000,
	"yearlyInterestRate": "5",
	"monthlyPayment": 1500000,
	"escrowPayment": 10000,
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| startingPrincipal | integer | Yes | Starting principal of the loan in cents. |
| yearlyInterestRate | string | Yes | Yearly interest rate of the loan (APR). Send as a percent ("5" would mean 5% yearly interest rate) |
| monthlyPayment | integer | Yes | Monthly payments made to the loan, in cents. |
| escrowPayment | integer | Yes | Additional, fixed payments that you make every month. For example, insurance, that are part of the monthly payments. |
| startDate | string | Yes | Start date of the loan. |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "durationMonths": 7,
    "totalExpenditure": 234054,
    "totalPaid": 10234054,
    "costOfCreditPercent": "1.02",
    "plan": [
        {
            "date": "1970-02-01T00:00:00Z",
            "payment": 1500000,
            "interest": 41667,
            "otherPayments": 10000,
            "paydown": 1448333,
            "principal": 8551667
        },
        {
            "date": "1970-03-01T00:00:00Z",
            "payment": 1500000,
            "interest": 35632,
            "otherPayments": 10000,
            "paydown": 1454368,
            "principal": 7097299
        }
    ]
}
```



### Response Fields

| Field | Type | Description |
|---|---|---|
| durationMonths | integer | Calculated duration of the loan in months. |
| totalExpenditure | integer | Total money paid to non-principal payments at the end of the loan. This includes interest and escrow payments. |
| totalPaid | integer | Total money paid at the end of the loan, including interest, principal and escrow payments. |
| costOfCreditPercent | string | How much more was paid than what was loaned. Calculated as the total paid over the loan's starting principal. For example "5" means a 5% cost of credit. |
| plan | array of plan statuses | An array of monthly statuses of the Payment plan. |

**Plan Status Fields**

| Field | Type | Description |
|---|---|---|
| date | string | The date of this monthly status. ISO 8601. |
| payment | integer | The payment made this month. |
| interest | integer | The interest accrued this month. |
| otherPayments | string | The amount that went into other payments this month (escrow payments like insurance). |
| paydown | array of plan statuses | The amount that was paid to principal this month. |
| principal | array of plan statuses | The loan principal at the end of this month. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST /app/savings/save`

Calculate and then save a Savings Plan to database.

---

### URL

```http
POST /app/savings/save
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
    "name": "Test",
	"startingCapital": 700000,
	"yearlyInterestRate": "4.75",
    "interestRateType": "APR",
	"monthlyContribution": 15000,
	"durationYears": 1,
    "taxRate": "5",
    "yearlyInflationRate": "6",
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| Name | string | Yes | The name of the savings plan. Doesn't need to be unique. |
| startingCapital | integer | Yes | How much money was deposited at the start of the term, in cents. For example, startingCapital = 100 would mean $1. |
| yearlyInterestRate | string | Yes | The yearly interest rate for the savings plan. Send as a percent. For example, "6.25" would be a 6.25% interest rate. |
| startDate | integer | Yes | The start date of the savings plan. |
| durationYears | integer | Yes | The term you want to calculate in years. 1 means "calculate the savings plan for 1 year".  ISO 8601|
| interestRateType | string | No | Send "APR" or "APY", depending on the type of interest rate. If empty, it defaults to APY. |
| monthlyContribution | integer | No | The monthly deposits that will be made (if any). Defaults to 0 if not in the request. The amount is in cents (e.g. 15000 means $150) |
| taxRate | string | No | The tax rate paid on returns. Send as a percent. For example, "5" means a 5% tax rate. Defaults to 0% if not in the request. |
| yearlyInflationRate | string | No | The yearly inflation rate, used for rate of return calculations. Send as a percent. For example, "6" means a 6% yearly inflation rate. Defaults to 0% if not in the request..  |

---

### Response

### Success Response

**Status Code:** `201 Created`

```json
{
    "id": "bd790a03-aabe-4a9d-861e-148fcd5adb46",
    "name": "test",
    "startingCapital": 700000,
    "yearlyInterestRate": "4.75",
    "interestRateType": "APY",
    "monthlyContribution": 15000,
    "durationYears": 1,
    "taxRate": "0",
    "yearlyInflationRate": "0",
    "startDate": "2026-01-31T18:00:00-06:00",
    "monthlyInterestRate": "0.3874684992129274856453709025059820372",
    "totalDeposited": 880000,
    "totalEarnings": 37136,
    "rateOfReturn": "4.22",
    "inflationAdjustedROR": "4.22"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the created Savings Plan. UUID v4. |
| name | string | The name of the created loan. |
| startingCapital | int | Confirmation of the starting capital saved in the database |
| yearlyInterestRate | string | Confirmation of the yearly interest rate saved in the database |
| interestRateType | string | Confirmation of the interest rate type saved in the database. |
| monthlyContribution | integer | Confirmation of the monthly contribution saved in the database. |
| durationYears | integer | Confirmation of the duration saved in the database. |
| taxRate | string | Confirmation of the tax rate saved in the database. |
| yearlyInflationRate | string | Confirmation of the yearly inflation rate saved in the database. |
| startDate | string | Confirmation of the start date saved in the database. ISO 8601 timestamp. |
| monthlyInterestRate | string | The monthly interest rate calculated from the yearly one, used for calculations. |
| totalDeposited | integer | The total money deposited on the account in cents. |
| totalEarnings | string | The total money earned in interest after the term. |
| rateOfReturn | string | The total savings in the account divided by what was actually deposited. Represents a percent (4.22 means 4.22% rate of return) |
| inflationAdjustedROR | string | The rate of retrun divided by the inflation at the end of the term. Represents a percent (4.22 means 4.22% rate of return). |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message dependimg on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST /app/loans/save`

Calculate and then save a Loan Payment Plan in database.

---

### URL

```http
POST /app/loans/save
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
    "name": "Test",
	"startingPrincipal": 10000000,
	"yearlyInterestRate": "5",
	"monthlyPayment": 1500000,
	"escrowPayment": 10000,
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| name | string | Yes | Name for the new Loan Payment Plan. |
| startingPrincipal | integer | Yes | Starting principal of the loan in cents. |
| yearlyInterestRate | string | Yes | Yearly interest rate of the loan (APR). Send as a percent ("5" would mean 5% yearly interest rate) |
| monthlyPayment | integer | Yes | Monthly payments made to the loan, in cents. |
| escrowPayment | integer | Yes | Additional, fixed payments that you make every month. For example, insurance, that are part of the monthly payments. |
| startDate | string | Yes | Start date of the loan. |

---

### Response

### Success Response

**Status Code:** `201 Created`

```json
{
    "id": "7632d773-4846-4b47-b3b9-3f3f06330072",
    "name": "Test",
    "startingPrincipal": 10000000,
    "yearlyInterestRate": "5",
    "monthlyPayment": 900076,
    "escrowPayment": 10000,
    "startDate": "1969-12-31T18:00:00-06:00",
    "durationMonths": 12,
    "totalExpenditure": 383416,
    "totalPaid": 10383416,
    "costOfCreditPercent": "3.83416398261762"
}
```



### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the created Loan Payment Plan. UUID v4. |
| name | string | The name of the created Loan Payment Plan. |
| startingPrincipal | integer | Confirmation of the starting principal saved on database. |
| yearlyInterestRate | integer | Confirmation of the yearly interest rate saved on database. |
| monthlyPayment | integer | Confirmation of the monthly payment saved on database. |
| escrowPayment | integer | Confirmation of the escrow payment saved on database. |
| startDate | integer | Confirmation of the start date saved on database. ISO 8601 timestamp. |
| durationMonths | integer | Calculated duration of the loan in months. |
| totalExpenditure | integer | The total non-principal expenditures at the end of the loan. Includes interest and escrow payments. |
| totalPaid | integer | The toal money paid on the loan. Includes interest, escrow and principal payments. |
| costOfCreditPercent | string | Calculated cost of credit as a percent. For example, "3.33" means a 3.33% cost of credit. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```

## Collaborators</h2>
<table>
  <tr>
    <td align="center">
      <a href="#">
        <img src="https://avatars.githubusercontent.com/u/35672719?s=48&v=4" width="100px;" alt="Rafael Mazariegos picture"/><br>
        <sub>
          <b>Rafael Mazariegos (Mr-Rafael)</b>
        </sub>
      </a>
    </td>
  </tr>
</table>