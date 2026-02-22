import { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import LoadingSpinner from '../ui/LoadingSpinner'

interface ProtectedRouteProps {
  children: ReactNode
  adminAllowed?: boolean
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, adminAllowed = false }) => {
  const { isAuthenticated, isAdmin, isLoading } = useAuth()
  const location = useLocation()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-gray-600">Checking authentication...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    // Redirect to login page with return url
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  // If user is admin and this route doesn't allow admins, redirect to admin dashboard
  if (isAdmin && !adminAllowed) {
    return <Navigate to="/admin/dashboard" replace />
  }

  // If user is regular user and trying to access admin-only route
  if (!isAdmin && adminAllowed) {
    return <Navigate to="/home" replace />
  }

  return <>{children}</>
}

export default ProtectedRoute

