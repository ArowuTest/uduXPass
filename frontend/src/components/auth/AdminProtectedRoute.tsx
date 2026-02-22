import { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import LoadingSpinner from '../ui/LoadingSpinner'

interface AdminProtectedRouteProps {
  children: ReactNode
  requiredPermissions?: string[]
  requiredRoles?: string[]
}

const AdminProtectedRoute: React.FC<AdminProtectedRouteProps> = ({ 
  children, 
  requiredPermissions = [], 
  requiredRoles = [] 
}) => {
  const { isAuthenticated, isAdmin, isLoading, canAccess } = useAuth()
  const location = useLocation()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  // Not authenticated at all - redirect to admin login
  if (!isAuthenticated) {
    return <Navigate to="/admin/login" state={{ from: location }} replace />
  }

  // Authenticated but not an admin - redirect to user home
  if (!isAdmin) {
    return <Navigate to="/" replace />
  }

  // Admin but doesn't have required permissions/roles
  if (requiredPermissions.length > 0 || requiredRoles.length > 0) {
    if (!canAccess(requiredPermissions, requiredRoles)) {
      return <Navigate to="/admin/dashboard" replace />
    }
  }

  return <>{children}</>
}

export default AdminProtectedRoute

