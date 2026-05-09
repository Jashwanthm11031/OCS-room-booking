import { useState, useEffect } from 'react'
import Navbar from '../../components/Navbar'
import { getAllBookings, adminCancelBooking } from '../../api/bookings'

export default function ManageBookings() {
  const [bookings, setBookings] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [filter, setFilter] = useState('')

  const fetchBookings = async () => {
    try {
      const res = await getAllBookings()
      setBookings(res.data.data || [])
    } catch { setError('Failed to load bookings') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchBookings() }, [])

  const handleCancel = async (id) => {
    if (!confirm('Cancel this booking?')) return
    try {
      await adminCancelBooking(id)
      fetchBookings()
    } catch (err) { setError(err.response?.data?.error || 'Failed to cancel booking') }
  }

  const filtered = filter
    ? bookings.filter(b =>
        b.room?.room_name?.toLowerCase().includes(filter.toLowerCase()) ||
        b.user?.name?.toLowerCase().includes(filter.toLowerCase()) ||
        b.purpose?.toLowerCase().includes(filter.toLowerCase()) ||
        b.date?.includes(filter)
      )
    : bookings

  const confirmed = filtered.filter(b => b.status === 'confirmed').length
  const cancelled = filtered.filter(b => b.status === 'cancelled').length

  return (
    <>
      <Navbar />
      <div className="page-container">
        <h2>All Bookings</h2>

        <div className="stats-row">
          <div className="stat-card">
            <div className="stat-number">{filtered.length}</div>
            <div className="stat-label">Total</div>
          </div>
          <div className="stat-card stat-green">
            <div className="stat-number">{confirmed}</div>
            <div className="stat-label">Confirmed</div>
          </div>
          <div className="stat-card stat-red">
            <div className="stat-number">{cancelled}</div>
            <div className="stat-label">Cancelled</div>
          </div>
        </div>

        <div className="search-bar">
          <input
            type="text"
            placeholder="Filter by room, user, purpose, or date..."
            value={filter}
            onChange={e => setFilter(e.target.value)}
          />
        </div>

        {error && <div className="alert alert-error">{error}</div>}

        {loading ? (
          <div className="loading">Loading bookings...</div>
        ) : filtered.length === 0 ? (
          <div className="empty-state">No bookings found.</div>
        ) : (
          <div className="table-wrapper">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Room</th>
                  <th>Block</th>
                  <th>Booked By</th>
                  <th>Date</th>
                  <th>Time</th>
                  <th>Purpose</th>
                  <th>Participants</th>
                  <th>Status</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map(b => (
                  <tr key={b.id}>
                    <td>{b.room?.room_name || '—'}</td>
                    <td>{b.room?.block?.name || '—'}</td>
                    <td>
                      <div>{b.user?.name || '—'}</div>
                      <div className="sub-text">{b.user?.email || ''}</div>
                    </td>
                    <td>{b.date}</td>
                    <td>{b.start_time?.slice(0, 5)} – {b.end_time?.slice(0, 5)}</td>
                    <td><span className="purpose-badge">{b.purpose}</span></td>
                    <td>{b.participant_count}</td>
                    <td>
                      <span className={`status-badge ${b.status}`}>{b.status}</span>
                    </td>
                    <td>
                      {b.status === 'confirmed' && (
                        <button className="btn btn-danger btn-sm" onClick={() => handleCancel(b.id)}>
                          Cancel
                        </button>
                      )}
                    </td>
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
