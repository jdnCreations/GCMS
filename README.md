# BOARD GAME CAFE RESERVATION SYSTEM

A web application for managing board game reservations at cafes. Users can browse games and make reservations online.

## Features

- User registration and authentication
- Make and manage reservations
- Admin dashboard

## Tech Stack

### Frontend

- Next.js
- Typescript
- Tailwind CSS
- React
- Axios

### Backend

- Go
- PostgreSQL

## Prerequisites

### Frontend Requirements

- Node.js (v22.3.0 or higher)
- npm (v10.8.1 or higher)
- TypeScript (v5.7.2 or higher)
- Next.js (v15.1.3)
- Tailwind CSS (v3.4.17)

### Backend Requirements

- Go (v1.23.3 or higher)
- PostgreSQL (v17.1 or higher)

### Development Tools

- Git
- A code editor (VS Code recommended)

## Installation & Setup

### Backend Setup

1. Clone the repository

```bash
git clone https://github.com/jdnCreations/GCMS
cd GCMS/backend
```

2. Install dependencies

```bash
go mod download
```

3. Create a .env file in backend directory

```
DATABASE_URL=postgres://<user>:<password>@localhost:5432/<dbname>?sslmode=disable
PORT=8080
TIMEZONE=<yourtimezone>
SECRET=<secret>
```

### Frontend Setup

1. Navigate to frontend directory

```bash
cd ../frontend
```

2. Install dependencies

```bash
npm install
```

3. Create a .env.local file

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Database Setup

1. Start the PostgreSQL service

```bash
# On linux/WSL
sudo service postgresql start

# On macOS (if installed via Homebrew)
brew services start postgresql

# On Windows (through Task Manager > Services)
# Find postgresql-x64-14 and click "Start"
```

2. Create a new database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE name_db;

# Verify database creation
\l
```

3. Run migrations using Goose

```bash
# must be in backend folder
goose -dir sql/schema postgres <your postgres url> up

# Check migration status
goose -dir sql/schema postgres <your postgres url> status
```

## Running the Application

### Start the Backend

```bash
cd backend
go run .
```

### Start the frontend

```bash
cd frontend
npm run dev
```
