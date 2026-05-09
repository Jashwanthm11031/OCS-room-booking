import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function ProtectedRoute({ children, adminOnly = false, coreOrAdmin = false }) {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (adminOnly && user?.role !== 'admin') {
    return <Navigate to="/dashboard" replace />
  }

  if (coreOrAdmin && user?.role !== 'admin' && user?.role !== 'core') {
    return <Navigate to="/dashboard" replace />
  }

  return children
}
