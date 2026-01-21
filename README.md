# Simple Procurement System

Technical Test for Junior Fullstack Engineer.

## Prerequisites
- Go 1.18+
- PostgreSQL
- `git`

## Setup

1. **Database**
   Create a database named `procurement_db` in your PostgreSQL.

2. **Environment Variables**
   Modify `.env` file if your PostgreSQL credentials are different:
   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=procurement_db
   DB_PORT=5432
   JWT_SECRET=supersecretkey
   PORT=3000
   WEBHOOK_URL=https://webhook.site/your-id
   ```

3. **Backend**
   ```bash
   go run cmd/api/main.go
   ```

4. **Frontend**
   Open `public/index.html` in your browser. Or serve it using a simple server:
   ```bash
   # example using python
   python -m http.server 8081 --directory public
   ```

## Features Implemented
- **Backend (Go/Fiber/GORM)**:
    - JWT Authentication (Register/Login).
    - Database Schema with UUID for Transactions.
    - Master Data CRUD (Suppliers & Items).
    - **ACID Transaction**: Create Purchase updates Item Stock and calculates totals server-side.
    - **Webhook Bonus**: Sends JSON notification after successful purchase.
- **Frontend (jQuery)**:
    - Dynamic SPA-style UI.
    - Token-based AJAX requests.
    - **Cart Logic**: Add/Remove items dynamically.
    - **Event Delegation**: Bonus requirement for dynamic elements.
    - **Reusable AJAX**: Bonus requirement for modular code.
    - **SweetAlert2**: Bonus requirement for better UX.

## User Account (Seed data)
Register a new user via the interface to start testing.
Need to add initial items/suppliers? Use the API or just register and you might need a way to seed data.
I've added simple forms for Register/Login. To add Items/Suppliers without a dedicated UI page (not required by the prompt but helpful), you can use Postman or I can add a quick seed script.
Actually, I'll add a simple "Admin" way to add items in the next step or just expect it to be done via API.
Wait, the user wants "CRUD for Items and Suppliers" in Backend. I implemented it. I should probably add a way to add them in UI too for completeness.
