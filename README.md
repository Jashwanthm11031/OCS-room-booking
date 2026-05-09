# OCS IITH Room Booking System

A full-stack room booking system for the Office of Career Services at IIT Hyderabad.

## Tech Stack

| Layer    | Technology                   |
|----------|------------------------------|
| Backend  | Go 1.21, Gin, GORM, JWT      |
| Database | PostgreSQL 14+               |
| Frontend | React 18, Vite, React Router |

---

## Roles

| Role    | Permissions                                               |
|---------|-----------------------------------------------------------|
| admin   | Full access: manage users, rooms, view/cancel all bookings |
| core    | Search rooms, create bookings, cancel own bookings        |
| viewer  | Read-only view of all bookings                            |

---

## Project Structure

```
ocs-room-booking/
├── backend/
│   ├── config/        # Env config loader
│   ├── db/            # GORM PostgreSQL connection
│   ├── handlers/      # HTTP handlers (auth, users, rooms, bookings)
│   ├── middleware/     # JWT auth + role guard
│   ├── migrations/    # Reference SQL (migrations run inline in main.go)
│   ├── models/        # GORM models
│   ├── repository/    # DB query layer
│   ├── services/      # Business logic
│   ├── utils/         # JWT, bcrypt, response helpers
│   ├── main.go        # Entry point, migrations, seeding, routing
│   ├── .env.example   # Copy to .env before running
│   └── go.mod
└── frontend/
    ├── src/
    │   ├── api/       # Axios API calls
    │   ├── components/# Navbar, RoomCard, BookingForm, ProtectedRoute
    │   ├── context/   # AuthContext (JWT + user state)
    │   ├── pages/     # Login, Dashboard, RoomSearch, BookingHistory
    │   │   └── admin/ # ManageUsers, ManageRooms, ManageBookings
    │   └── styles/    # main.css
    ├── App.jsx
    ├── main.jsx
    ├── index.html
    └── package.json
```

---

## Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### 1. Database

```sql
CREATE DATABASE ocs_booking;
```

### 2. Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your DB credentials and a strong JWT_SECRET

go mod tidy
go run main.go
# Server starts on http://localhost:8080
# Tables are created and seed data (admin + all IITH rooms) inserted automatically
```

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
# App starts on http://localhost:5173
```

---

## API Endpoints

### Public
| Method | Path             | Description   |
|--------|------------------|---------------|
| POST   | /api/auth/login  | Login, get JWT |

### Authenticated (all roles)
| Method | Path                  | Description                  |
|--------|-----------------------|------------------------------|
| GET    | /api/rooms/search     | Search available rooms       |
| GET    | /api/rooms/:id        | Get single room              |
| GET    | /api/blocks           | List all blocks              |
| GET    | /api/bookings/my      | My bookings (core)           |
| GET    | /api/bookings/all     | All bookings (viewer+admin)  |
| DELETE | /api/bookings/:id     | Cancel own booking (core)    |

### core + admin
| Method | Path          | Description    |
|--------|---------------|----------------|
| POST   | /api/bookings | Create booking |

### Admin only
| Method | Path                    | Description           |
|--------|-------------------------|-----------------------|
| GET    | /api/admin/users        | List all users        |
| POST   | /api/admin/users        | Create user           |
| PATCH  | /api/admin/users/:id    | Update user / toggle active |
| POST   | /api/admin/rooms        | Add room              |
| PATCH  | /api/admin/rooms/:id    | Update room           |
| DELETE | /api/admin/rooms/:id    | Delete room           |
| GET    | /api/admin/bookings     | All bookings          |
| DELETE | /api/admin/bookings/:id | Cancel any booking    |

---

## Seeded Data

On first run, `main.go` automatically creates:
- **1 admin user** (credentials from `.env`)
- **11 blocks**: A Block, B Block, C Block, CSE Block, LHC, BT/BM, CY, EE, MA, MSME, PH
- **54 rooms** across all blocks with real IITH capacities
