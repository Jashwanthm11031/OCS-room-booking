import { useAuth } from '../context/AuthContext'
import { Link } from 'react-router-dom'
import Navbar from '../components/Navbar'

export default function Dashboard() {
  const { user } = useAuth()

  return (
    <>
      <Navbar />
      <div className="page-container">
        <div className="dashboard-welcome">
          <h2>Welcome, {user?.name}!</h2>
          <p className="subtitle">IIT Hyderabad — Office of Career Services Room Booking</p>
        </div>

        <div className="dashboard-cards">
          {(user?.role === 'core' || user?.role === 'admin') && (
            <Link to="/rooms" className="dash-card">
              <div className="dash-card-icon">🔍</div>
              <h3>Search Rooms</h3>
              <p>Find and book available rooms across all IITH blocks</p>
            </Link>
          )}
          {user?.role !== 'admin' && (
            <Link to="/bookings" className="dash-card">
              <div className="dash-card-icon">📅</div>
              <h3>{user?.role === 'viewer' ? 'All Bookings' : 'My Bookings'}</h3>
              <p>View booking history and manage your reservations</p>
            </Link>
          )}
          {user?.role === 'admin' && (
            <>
              <Link to="/admin/users" className="dash-card">
                <div className="dash-card-icon">👥</div>
                <h3>Manage Users</h3>
                <p>Create and manage user accounts</p>
              </Link>
              <Link to="/admin/rooms" className="dash-card">
                <div className="dash-card-icon">🏫</div>
                <h3>Manage Rooms</h3>
                <p>Add, edit, or remove rooms across all blocks</p>
              </Link>
              <Link to="/admin/bookings" className="dash-card">
                <div className="dash-card-icon">📋</div>
                <h3>All Bookings</h3>
                <p>View and manage all bookings system-wide</p>
              </Link>
            </>
          )}
        </div>
      </div>
    </>
  )
}
