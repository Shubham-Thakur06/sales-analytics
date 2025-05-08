# sales-analytics

A Go-based REST API for analyzing sales data with support for batch processing and revenue analytics.

## Features

- Batch processing of CSV data with configurable batch size
- Automated data refresh using cron jobs
- Revenue analytics by:
  - Total revenue
  - Product-wise breakdown
  - Category-wise breakdown
  - Regional breakdown
- PostgreSQL database with GORM ORM
- Configurable through environment variables

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

## Configuration

Create a `.env` file in the project root with the following configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=sales_insights
DB_SSL_MODE=disable

# Application Configuration
APP_PORT=8080
CSV_FILE_PATH=path/to/data.csv
REFRESH_CRON=0 0 * * * # Run at midnight every day

# Processing Configuration
BATCH_SIZE=1000 # Number of records to process in each batch
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/Shubham-Thakur06/sales-analytics.git
cd sales-analytics
```

2. Install dependencies:
```bash
go mod download
```

3. Create the PostgreSQL database:
```sql
CREATE DATABASE sales_insights;
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/refresh` | Triggers manual refresh of CSV data |
| GET | `/api/v1/revenue` | Get total revenue for date range |
| GET | `/api/v1/revenue/product` | Get revenue breakdown by product |
| GET | `/api/v1/revenue/category` | Get revenue breakdown by category |

All revenue endpoints accept query parameters:
- `start_date`: Start date (YYYY-MM-DD)
- `end_date`: End date (YYYY-MM-DD)

### Data Refresh

- **POST** `/api/v1/refresh`
  - Triggers a manual refresh of the data from CSV
  - Response:
    ```json
    {
      "message": "Data refresh started"
    }
    ```

### Revenue Analytics

All revenue endpoints accept date range parameters:
- `start_date`: Start date in YYYY-MM-DD format
- `end_date`: End date in YYYY-MM-DD format

1. **GET** `/api/v1/revenue`
   - Get total revenue for the specified date range
   - Response:
     ```json
     {
       "total_revenue": 150000.50,
       "start_date": "2023-01-01T00:00:00Z",
       "end_date": "2023-12-31T00:00:00Z"
     }
     ```

2. **GET** `/api/v1/revenue/product`
   - Get revenue breakdown by product
   - Response:
     ```json
     [
       {
         "product_id": "P123",
         "product_name": "Product A",
         "revenue": 50000.25
       }
     ]
     ```

3. **GET** `/api/v1/revenue/category`
   - Get revenue breakdown by product category
   - Response:
     ```json
     [
       {
         "category": "Electronics",
         "revenue": 75000.00
       }
     ]
     ```

4. **GET** `/api/v1/revenue/region`
   - Get revenue breakdown by region
   - Response:
     ```json
     [
       {
         "region": "North America",
         "revenue": 100000.75
       }
     ]
     ```

### Error Responses

The API returns appropriate HTTP status codes and error messages:

- 400 Bad Request: For invalid input (e.g., invalid date format)
  ```json
  {
    "error": "Invalid start date '2023-13-45'. Date must be in format YYYY-MM-DD"
  }
  ```

- 500 Internal Server Error: For server-side errors
  ```json
  {
    "error": "Failed to calculate total revenue"
  }
  ```

## CSV Data Format

The CSV file should contain the following columns:
1. Order ID
2. Product ID
3. Customer ID
4. Product Name
5. Category
6. Region
7. Date of Sale (YYYY-MM-DD)
8. Quantity
9. Unit Price
10. Discount
11. Shipping Cost
12. Customer Name
13. Customer Email
14. Customer Address

## Development

### Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/        # HTTP handlers
│   │   └── routes.go        # Route definitions
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── container/
│   │   └── container.go     # Dependency injection
│   ├── models/
│   │   ├── customer.go      # Data models
│   │   ├── order.go
│   │   ├── product.go
│   │   └── revenue.go       # Response models
│   └── services/
│       ├── loader.go        # CSV data loading
│       └── revenue.go       # Revenue calculations
├── .env.example             # Example configuration
├── go.mod                   # Go module file
└── README.md               # This file
```

## License

MIT License - See LICENSE file for details 