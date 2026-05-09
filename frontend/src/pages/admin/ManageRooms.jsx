import { useState, useEffect } from 'react'
import Navbar from '../../components/Navbar'
import { getAllRooms, getBlocks, createRoom, updateRoom, deleteRoom } from '../../api/rooms'

const PURPOSES = ['OA', 'Interview', 'PPT']

export default function ManageRooms() {
  const [rooms, setRooms] = useState([])
  const [blocks, setBlocks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [editRoom, setEditRoom] = useState(null)
  const [form, setForm] = useState({
    block_id: '', room_name: '', capacity: '', allowed_purposes: ['OA', 'Interview', 'PPT'], notes: ''
  })
  const [creating, setCreating] = useState(false)

  const fetchAll = async () => {
    try {
      const [roomsRes, blocksRes] = await Promise.all([getAllRooms(), getBlocks()])
      setRooms(roomsRes.data.data || [])
      setBlocks(blocksRes.data.data || [])
    } catch { setError('Failed to load data') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchAll() }, [])

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const togglePurpose = (p) => {
    const current = form.allowed_purposes
    if (current.includes(p)) {
      setForm({ ...form, allowed_purposes: current.filter(x => x !== p) })
    } else {
      setForm({ ...form, allowed_purposes: [...current, p] })
    }
  }

  const handleCreate = async (e) => {
    e.preventDefault()
    setError(''); setSuccess(''); setCreating(true)
    try {
      await createRoom({
        block_id: form.block_id,
        room_name: form.room_name,
        capacity: parseInt(form.capacity),
        allowed_purposes: form.allowed_purposes,
        notes: form.notes,
      })
      setSuccess('Room created successfully')
      setForm({ block_id: '', room_name: '', capacity: '', allowed_purposes: ['OA', 'Interview', 'PPT'], notes: '' })
      fetchAll()
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create room')
    } finally { setCreating(false) }
  }

  const handleToggleAvailable = async (room) => {
    try {
      await updateRoom(room.id, { is_available: !room.is_available })
      fetchAll()
    } catch (err) { setError(err.response?.data?.error || 'Failed to update room') }
  }

  const handleDelete = async (id) => {
    if (!confirm('Delete this room? All its bookings will also be removed.')) return
    try {
      await deleteRoom(id)
      fetchAll()
    } catch (err) { setError(err.response?.data?.error || 'Failed to delete room') }
  }

  return (
    <>
      <Navbar />
      <div className="page-container">
        <h2>Manage Rooms</h2>

        <div className="form-card">
          <h3>Add New Room</h3>
          {error && <div className="alert alert-error">{error}</div>}
          {success && <div className="alert alert-success">{success}</div>}
          <form onSubmit={handleCreate}>
            <div className="form-grid">
              <div className="form-group">
                <label>Block</label>
                <select name="block_id" value={form.block_id} onChange={handleChange} required>
                  <option value="">Select Block</option>
                  {blocks.map(b => <option key={b.id} value={b.id}>{b.name}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Room Name</label>
                <input type="text" name="room_name" value={form.room_name} onChange={handleChange} required placeholder="e.g. LHC-01" />
              </div>
              <div className="form-group">
                <label>Capacity</label>
                <input type="number" name="capacity" value={form.capacity} onChange={handleChange} required min="1" placeholder="e.g. 80" />
              </div>
              <div className="form-group">
                <label>Notes</label>
                <input type="text" name="notes" value={form.notes} onChange={handleChange} placeholder="Optional notes" />
              </div>
            </div>
            <div className="form-group">
              <label>Allowed Purposes</label>
              <div className="checkbox-group">
                {PURPOSES.map(p => (
                  <label key={p} className="checkbox-label">
                    <input
                      type="checkbox"
                      checked={form.allowed_purposes.includes(p)}
                      onChange={() => togglePurpose(p)}
                    />
                    {p}
                  </label>
                ))}
              </div>
            </div>
            <button type="submit" className="btn btn-primary" disabled={creating}>
              {creating ? 'Adding...' : 'Add Room'}
            </button>
          </form>
        </div>

        <div className="section-header"><h3>All Rooms ({rooms.length})</h3></div>
        {loading ? (
          <div className="loading">Loading rooms...</div>
        ) : (
          <div className="table-wrapper">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Room</th>
                  <th>Block</th>
                  <th>Capacity</th>
                  <th>Purposes</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {rooms.map(r => (
                  <tr key={r.id}>
                    <td><strong>{r.room_name}</strong></td>
                    <td>{r.block?.name}</td>
                    <td>{r.capacity}</td>
                    <td>{r.allowed_purposes?.join(', ')}</td>
                    <td>
                      <span className={`status-badge ${r.is_available ? 'confirmed' : 'cancelled'}`}>
                        {r.is_available ? 'Available' : 'Unavailable'}
                      </span>
                    </td>
                    <td className="action-cell">
                      <button
                        className={`btn btn-sm ${r.is_available ? 'btn-warning' : 'btn-success'}`}
                        onClick={() => handleToggleAvailable(r)}
                      >
                        {r.is_available ? 'Disable' : 'Enable'}
                      </button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDelete(r.id)}>
                        Delete
                      </button>
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
