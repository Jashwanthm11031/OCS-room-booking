import { useNavigate } from 'react-router-dom'

export default function RoomCard({ room, searchParams }) {
  const navigate = useNavigate()

  const handleBook = () => {
    navigate('/book', { state: { room, searchParams } })
  }

  return (
    <div className="room-card">
      <div className="room-card-header">
        <h3>{room.room_name}</h3>
        <span className="block-badge">{room.block?.name}</span>
      </div>
      <div className="room-card-body">
        <div className="room-info">
          <span>👥 Capacity: <strong>{room.capacity}</strong></span>
          <span>📋 Purposes: <strong>{room.allowed_purposes?.join(', ')}</strong></span>
          {room.notes && <span>📝 {room.notes}</span>}
        </div>
      </div>
      <div className="room-card-footer">
        <span className={`status-badge ${room.is_available ? 'available' : 'unavailable'}`}>
          {room.is_available ? 'Available' : 'Unavailable'}
        </span>
        {room.is_available && (
          <button className="btn btn-primary" onClick={handleBook}>
            Book Room
          </button>
        )}
      </div>
    </div>
  )
}
