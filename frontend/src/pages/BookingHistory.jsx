import { useState, useEffect } from 'react'
import Navbar from '../components/Navbar'
import { getMyBookings, cancelMyBooking, getAllBookings } from '../api/bookings'
import { useAuth } from '../context/AuthContext'

export default function BookingHistory() {
  const { user } = useAuth()
  const isViewer = user?.role === 'viewer'

  const [bookings, setBookings] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const fetchBookings = async () => {
    setLoading(true)
    try {
      // Viewers see all bookings (read-only); core members see only their own
      const res = isViewer ? await getAllBookings() : await getMyBookings()
      setBookings(res.data.data || [])
    } catch {
      setError('Failed to load bookings')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchBookings() }, [])

  const handleCancel = async (id) => {
    if (!confirm('Cancel this booking?')) return
    try {
      await cancelMyBooking(id)
      fetchBookings()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to cancel booking')
    }
  }

  return (
    <>
      <Navbar />
      <div className="page-container">
        <h2>{isViewer ? 'All Bookings (View Only)' : 'My Bookings'}</h2>
        {error && <div className="alert alert-error">{error}</div>}
        {loading ? (
          <div className="loading">Loading bookings...</div>
        ) : bookings.length === 0 ? (
          <div className="empty-state">
            No bookings found.{' '}
            {!isViewer && <a href="/rooms">Search for a room</a>}
          </div>
        ) : (
          <div className="table-wrapper">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Room</th>
                  <th>Block</th>
                  {isViewer && <th>Booked By</th>}
                  <th>Date</th>
                  <th>Time</th>
                  <th>Purpose</th>
                  <th>Participants</th>
                  <th>Status</th>
                  {!isViewer && <th>Action</th>}
                </tr>
              </thead>
              <tbody>
                {bookings.map(b => (
                  <tr key={b.id}>
                    <td>{b.room?.room_name || '—'}</td>
                    <td>{b.room?.block?.name || '—'}</td>
                    {isViewer && (
                      <td>
                        <div>{b.user?.name || '—'}</div>
                        <div className="sub-text">{b.user?.email || ''}</div>
                      </td>
                    )}
                    <td>{b.date}</td>
                    <td>{b.start_time?.slice(0, 5)} – {b.end_time?.slice(0, 5)}</td>
                    <td><span className="purpose-badge">{b.purpose}</span></td>
                    <td>{b.participant_count}</td>
                    <td>
                      <span className={`status-badge ${b.status}`}>{b.status}</span>
                    </td>
                    {!isViewer && (
                      <td>
                        {b.status === 'confirmed' && (
                          <button
                            className="btn btn-danger btn-sm"
                            onClick={() => handleCancel(b.id)}
                          >
                            Cancel
                          </button>
                        )}
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  )
}
