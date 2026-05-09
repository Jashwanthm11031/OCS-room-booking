import api from './index'

export const createBooking = (data) =>
  api.post('/bookings', data)

export const getMyBookings = () =>
  api.get('/bookings/my')

export const cancelMyBooking = (id) =>
  api.delete(`/bookings/${id}`)

// Viewer + Admin: see all bookings
export const getAllBookings = () =>
  api.get('/bookings/all')

// Admin only
export const adminCancelBooking = (id) =>
  api.delete(`/admin/bookings/${id}`)
