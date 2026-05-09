import { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { createBooking } from '../api/bookings'

export default function BookingForm() {
  const { state } = useLocation()
  const navigate = useNavigate()
  const room = state?.room
  const pre = state?.searchParams || {}

  const [form, setForm] = useState({
    date: pre.date || '',
    start_time: pre.start_time || '',
    end_time: pre.end_time || '',
    purpose: pre.purpose || 'OA',
    participant_count: '',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  if (!room) {
    navigate('/rooms')
    return null
  }

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await createBooking({
        room_id: room.id,
        date: form.date,
        start_time: form.start_time,
        end_time: form.end_time,
        purpose: form.purpose,
        participant_count: parseInt(form.participant_count),
      })
      navigate('/bookings', { state: { success: true } })
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create booking')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page-container">
      <div className="form-card">
        <h2>Book Room</h2>
        <div className="room-summary">
          <strong>{room.room_name}</strong> — {room.block?.name} (Capacity: {room.capacity})
        </div>

        {error && <div className="alert alert-error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Date</label>
            <input type="date" name="date" value={form.date} onChange={handleChange} required
              min={new Date().toISOString().split('T')[0]} />
          </div>
          <div className="form-row">
            <div className="form-group">
              <label>Start Time</label>
              <input type="time" name="start_time" value={form.start_time} onChange={handleChange} required />
            </div>
            <div className="form-group">
              <label>End Time</label>
              <input type="time" name="end_time" value={form.end_time} onChange={handleChange} required />
            </div>
          </div>
          <div className="form-group">
            <label>Purpose</label>
            <select name="purpose" value={form.purpose} onChange={handleChange} required>
              {room.allowed_purposes?.map(p => (
                <option key={p} value={p}>{p}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Number of Participants (max: {room.capacity})</label>
            <input type="number" name="participant_count" value={form.participant_count}
              onChange={handleChange} required min="1" max={room.capacity} />
          </div>
          <div className="form-actions">
            <button type="button" className="btn btn-outline" onClick={() => navigate(-1)}>Cancel</button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Booking...' : 'Confirm Booking'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
