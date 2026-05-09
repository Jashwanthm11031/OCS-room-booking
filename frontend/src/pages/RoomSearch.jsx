import { useState, useEffect } from 'react'
import Navbar from '../components/Navbar'
import RoomCard from '../components/RoomCard'
import { searchRooms, getBlocks } from '../api/rooms'

export default function RoomSearch() {
  const [blocks, setBlocks] = useState([])
  const [rooms, setRooms] = useState([])
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)
  const [filters, setFilters] = useState({
    block: '',
    capacity: '',
    purpose: '',
    date: '',
    start_time: '',
    end_time: '',
  })

  useEffect(() => {
    getBlocks().then(r => setBlocks(r.data.data || []))
  }, [])

  const handleChange = (e) => {
    setFilters({ ...filters, [e.target.name]: e.target.value })
  }

  const handleSearch = async (e) => {
    e.preventDefault()
    setLoading(true)
    setSearched(true)
    try {
      const params = {}
      if (filters.block) params.block = filters.block
      if (filters.capacity) params.capacity = filters.capacity
      if (filters.purpose) params.purpose = filters.purpose
      if (filters.date) params.date = filters.date
      if (filters.start_time) params.start_time = filters.start_time
      if (filters.end_time) params.end_time = filters.end_time
      const res = await searchRooms(params)
      setRooms(res.data.data || [])
    } catch {
      setRooms([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Navbar />
      <div className="page-container">
        <h2>Search Rooms</h2>

        <div className="search-form-card">
          <form onSubmit={handleSearch}>
            <div className="form-grid">
              <div className="form-group">
                <label>Block</label>
                <select name="block" value={filters.block} onChange={handleChange}>
                  <option value="">All Blocks</option>
                  {blocks.map(b => (
                    <option key={b.id} value={b.id}>{b.name}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>Min Capacity</label>
                <input type="number" name="capacity" value={filters.capacity}
                  onChange={handleChange} placeholder="e.g. 50" min="1" />
              </div>
              <div className="form-group">
                <label>Purpose</label>
                <select name="purpose" value={filters.purpose} onChange={handleChange}>
                  <option value="">Any Purpose</option>
                  <option value="OA">Online Assessment (OA)</option>
                  <option value="Interview">Interview</option>
                  <option value="PPT">Pre-Placement Talk (PPT)</option>
                </select>
              </div>
              <div className="form-group">
                <label>Date</label>
                <input type="date" name="date" value={filters.date} onChange={handleChange}
                  min={new Date().toISOString().split('T')[0]} />
              </div>
              <div className="form-group">
                <label>Start Time</label>
                <input type="time" name="start_time" value={filters.start_time} onChange={handleChange} />
              </div>
              <div className="form-group">
                <label>End Time</label>
                <input type="time" name="end_time" value={filters.end_time} onChange={handleChange} />
              </div>
            </div>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Searching...' : 'Search Rooms'}
            </button>
          </form>
        </div>

        {searched && (
          <div className="search-results">
            {loading ? (
              <div className="loading">Searching rooms...</div>
            ) : rooms.length === 0 ? (
              <div className="empty-state">No rooms found matching your criteria.</div>
            ) : (
              <>
                <p className="results-count">{rooms.length} room(s) found</p>
                <div className="room-grid">
                  {rooms.map(room => (
                    <RoomCard key={room.id} room={room} searchParams={filters} />
                  ))}
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </>
  )
}
