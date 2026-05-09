import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import ProtectedRoute from './components/ProtectedRoute'

import Login          from './pages/Login'
import Dashboard      from './pages/Dashboard'
import RoomSearch     from './pages/RoomSearch'
import BookingHistory from './pages/BookingHistory'
import BookingForm    from './components/BookingForm'

import ManageUsers    from './pages/admin/ManageUsers'
import ManageRooms    from './pages/admin/ManageRooms'
import ManageBookings from './pages/admin/ManageBookings'

// Redirect based on role after login
function RootRedirect() {
  const { isAuthenticated, user } = useAuth()
  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (user?.role === 'admin') return <Navigate to="/admin/users" replace />
  return <Navigate to="/dashboard" replace />
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public */}
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<RootRedirect />} />

          {/* All authenticated roles */}
          <Route path="/dashboard" element={
            <ProtectedRoute><Dashboard /></ProtectedRoute>
          } />

          {/* Bookings: core sees own, viewer sees all (handled inside component) */}
          <Route path="/bookings" element={
            <ProtectedRoute><BookingHistory /></ProtectedRoute>
          } />

          {/* core + admin: search and book rooms */}
          <Route path="/rooms" element={
            <ProtectedRoute coreOrAdmin><RoomSearch /></ProtectedRoute>
          } />
          <Route path="/book" element={
            <ProtectedRoute coreOrAdmin><BookingForm /></ProtectedRoute>
          } />

          {/* Admin only */}
          <Route path="/admin/users" element={
            <ProtectedRoute adminOnly><ManageUsers /></ProtectedRoute>
          } />
          <Route path="/admin/rooms" element={
            <ProtectedRoute adminOnly><ManageRooms /></ProtectedRoute>
          } />
          <Route path="/admin/bookings" element={
            <ProtectedRoute adminOnly><ManageBookings /></ProtectedRoute>
          } />

          {/* Catch-all */}
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}
