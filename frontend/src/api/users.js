import api from './index'

export const getAllUsers = () =>
  api.get('/admin/users')

export const createUser = (data) =>
  api.post('/admin/users', data)

export const updateUser = (id, data) =>
  api.patch(`/admin/users/${id}`, data)
