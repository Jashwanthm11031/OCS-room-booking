import api from './index'

export const login = (email, password) =>
  api.post('/auth/login', { email, password })
