import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <Link to="/dashboard">OCS IITH Room Booking</Link>
      </div>
      <div className="navbar-links">
        {user?.role !== 'viewer' && (
          <>
            <Link to="/rooms">Search Rooms</Link>
            {user?.role !== 'admin' && <Link to="/bookings">My Bookings</Link>}
          </>
        )}
        {user?.role === 'viewer' && <Link to="/bookings">All Bookings</Link>}
        {user?.role === 'admin' && (
          <>
            <Link to="/admin/users">Users</Link>
            <Link to="/admin/rooms">Rooms</Link>
            <Link to="/admin/bookings">All Bookings</Link>
          </>
        )}
      </div>
      <div className="navbar-user">
        <span className="user-badge">{user?.name} ({user?.role})</span>
        <button onClick={handleLogout} className="btn btn-outline">Logout</button>
      </div>
    </nav>
  )
}
