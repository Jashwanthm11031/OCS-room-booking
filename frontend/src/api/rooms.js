import api from './index'

export const searchRooms = (params) =>
  api.get('/rooms/search', { params })

export const getRoom = (id) =>
  api.get(`/rooms/${id}`)

export const getBlocks = () =>
  api.get('/blocks')

// Admin: all rooms including unavailable
export const getAllRooms = () =>
  api.get('/admin/rooms')

// Admin CRUD
export const createRoom = (data) =>
  api.post('/admin/rooms', data)

export const updateRoom = (id, data) =>
  api.patch(`/admin/rooms/${id}`, data)

export const deleteRoom = (id) =>
  api.delete(`/admin/rooms/${id}`)
