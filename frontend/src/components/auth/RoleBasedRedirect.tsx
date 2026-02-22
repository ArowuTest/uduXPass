import { Navigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import LoadingSpinner from '../ui/LoadingSpinner'

const RoleBasedRedirect: React.FC = () => {
  const { isAuthenticated, isAdmin, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900">
        <div className="text-center">
          <div className="mb-8">
            <h1 className="text-6xl font-bold text-white mb-2">uduXPass</h1>
            <p className="text-xl text-purple-200">Premium Event Ticketing</p>
          </div>
          <LoadingSpinner size="lg" />
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    // Not authenticated - show public home page
    return <Navigate to="/home" replace />
  }

  if (isAdmin) {
    // Admin user - redirect to admin dashboard
    return <Navigate to="/admin/dashboard" replace />
  } else {
    // Regular user - redirect to user home
    return <Navigate to="/home" replace />
  }
}

export default RoleBasedRedirect

