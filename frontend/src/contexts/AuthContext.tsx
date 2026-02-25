import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authAPI, adminAuthAPI } from '../services/api'
import { User, AdminUser, AdminPermission, AdminRole } from '../types'

interface AuthContextType {
  user: User | null
  admin: AdminUser | null
  isAuthenticated: boolean
  isAdmin: boolean
  isLoading: boolean
  error: string | null
  login: (email: string, password: string) => Promise<void>
  register: (data: { email: string; password: string; firstName?: string; lastName?: string; phone?: string }) => Promise<{ success: boolean; error?: string }>
  logout: () => void
  adminLogin: (email: string, password: string) => Promise<void>
  adminLogout: () => void
  hasPermission: (permission: string) => boolean
  hasRole: (role: string) => boolean
  canAccess: (requiredPermissions?: string[], requiredRoles?: string[]) => boolean
  clearError: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

// Utility functions
const safeParseJSON = (jsonString: string, fallback: any = null): any => {
  try {
    return JSON.parse(jsonString)
  } catch {
    return fallback
  }
}

const isValidArray = (value: any): value is any[] => {
  return Array.isArray(value)
}

const validateAdminUser = (data: any): AdminUser | null => {
  if (!data || typeof data !== 'object') return null
  
  // Handle both camelCase and snake_case from backend
  const firstName = data.firstName || data.first_name
  const lastName = data.lastName || data.last_name
  
  if (!data.id || !data.email || !firstName || !lastName || !data.role) {
    return null
  }
  
  // FIXED: Grant all permissions to super_admin users
  let permissions = isValidArray(data.permissions) ? data.permissions : []
  
  if (data.role === 'super_admin') {
    // Super admin gets all permissions
    permissions = [
      'users_view', 'users_create', 'users_update', 'users_delete',
      'events_view', 'events_create', 'events_update', 'events_delete',
      'orders_view', 'orders_update', 'orders_delete',
      'tickets_view', 'tickets_update',
      'scanners_view', 'scanners_create', 'scanners_update', 'scanners_delete',
      'analytics_view', 'reports_view', 'settings_update',
      'admin_create', 'admin_update', 'admin_delete',
      // Also include dot notation versions for compatibility
      'users.view', 'users.create', 'users.update', 'users.delete',
      'events.view', 'events.create', 'events.update', 'events.delete',
      'orders.view', 'orders.update', 'orders.delete',
      'tickets.view', 'tickets.update',
      'scanners.view', 'scanners.create', 'scanners.update', 'scanners.delete',
      'analytics.view', 'reports.view', 'settings.update',
      'admin.create', 'admin.update', 'admin.delete'
    ]
  }
  
  return {
    id: data.id,
    email: data.email,
    firstName: firstName,
    lastName: lastName,
    role: data.role as AdminRole,
    permissions: permissions,
    isActive: Boolean(data.isActive || data.is_active),
    lastLoginAt: data.lastLoginAt || data.last_login || undefined,
    createdAt: data.createdAt || data.created_at || new Date().toISOString(),
    updatedAt: data.updatedAt || data.updated_at || new Date().toISOString()
  }
}

const validateUser = (data: any): User | null => {
  if (!data || typeof data !== 'object') return null
  
  const requiredFields = ['id', 'email']
  const hasRequiredFields = requiredFields.every(field => 
    data[field] && typeof data[field] === 'string'
  )
  
  if (!hasRequiredFields) return null
  
  return {
    id: data.id,
    email: data.email,
    phone: data.phone || undefined,
    firstName: data.firstName || undefined,
    lastName: data.lastName || undefined,
    isEmailVerified: Boolean(data.isEmailVerified),
    isPhoneVerified: Boolean(data.isPhoneVerified),
    role: data.role || 'user',
    createdAt: data.createdAt || new Date().toISOString(),
    updatedAt: data.updatedAt || new Date().toISOString()
  }
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [admin, setAdmin] = useState<AdminUser | null>(null)
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
  const [isAdmin, setIsAdmin] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    initializeAuth()
  }, [])

  const initializeAuth = async (): Promise<void> => {
    try {
      setIsLoading(true)
      setError(null)
      
      // Check for admin session first (higher priority)
      const adminToken = localStorage.getItem('adminToken')
      const adminDataString = localStorage.getItem('adminData')
      
      if (adminToken && adminDataString) {
        const adminData = safeParseJSON(adminDataString, null)
        const validatedAdmin = validateAdminUser(adminData)
        
        if (validatedAdmin) {
          setAdmin(validatedAdmin)
          setIsAdmin(true)
          setIsAuthenticated(true)
          setIsLoading(false)
          return
        } else {
          // Clean up invalid admin session
          localStorage.removeItem('adminToken')
          localStorage.removeItem('adminData')
        }
      }
      
      // Check for user session
      const userToken = localStorage.getItem('accessToken')
      const userDataString = localStorage.getItem('userData')
      
      if (userToken && userDataString) {
        const userData = safeParseJSON(userDataString, null)
        const validatedUser = validateUser(userData)
        
        if (validatedUser) {
          setUser(validatedUser)
          setIsAuthenticated(true)
          setIsAdmin(false)
          setIsLoading(false)
          return
        } else {
          // Clean up invalid user session
          localStorage.removeItem('accessToken')
          localStorage.removeItem('userData')
        }
      }
      
      // No valid session found
      setIsLoading(false)
    } catch (error) {
      console.error('Auth initialization error:', error)
      setError('Failed to initialize authentication')
      setIsLoading(false)
    }
  }

  const login = async (email: string, password: string): Promise<void> => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await authAPI.login({ email, password })
      
      if (response.success && response.data) {
        const { access_token, user: userData } = response.data
        const validatedUser = validateUser(userData)
        
        if (validatedUser && access_token) {
          localStorage.setItem('accessToken', access_token)
          localStorage.setItem('userData', JSON.stringify(validatedUser))
          
          setUser(validatedUser)
          setIsAuthenticated(true)
          setIsAdmin(false)
        } else {
          throw new Error('Invalid user data received')
        }
      } else {
        throw new Error(response.error || 'Login failed')
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Login failed'
      setError(errorMessage)
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  const register = async (data: {
    email: string;
    password: string;
    firstName?: string;
    lastName?: string;
    phone?: string;
  }): Promise<{ success: boolean; error?: string }> => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await authAPI.register(data)
      
      if (response.success && response.data) {
        const { access_token, user: userData } = response.data
        const validatedUser = validateUser(userData)
        
        if (validatedUser && access_token) {
          localStorage.setItem('accessToken', access_token)
          localStorage.setItem('userData', JSON.stringify(validatedUser))
          
          setUser(validatedUser)
          setIsAuthenticated(true)
          setIsAdmin(false)
          return { success: true }
        } else {
          const error = 'Invalid user data received'
          setError(error)
          return { success: false, error }
        }
      } else {
        const error = response.error || 'Registration failed'
        setError(error)
        return { success: false, error }
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Registration failed'
      setError(errorMessage)
      return { success: false, error: errorMessage }
    } finally {
      setIsLoading(false)
    }
  }

  const logout = (): void => {
    try {
      localStorage.removeItem('accessToken')
      localStorage.removeItem('userData')
      
      setUser(null)
      setIsAuthenticated(false)
      setIsAdmin(false)
      setError(null)
    } catch (error) {
      console.error('Logout error:', error)
    }
  }

  const adminLogin = async (email: string, password: string): Promise<void> => {
    console.log('[AuthContext] Starting admin login for:', email)
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await adminAuthAPI.login({ email, password })
      console.log('[AuthContext] Login API response:', response)
      
      if (response.success && response.data) {
        console.log('[AuthContext] Response data:', response.data)
        // Handle nested data structure from backend
        const actualData = response.data.data || response.data
        const { access_token, admin: adminData } = actualData
        console.log('[AuthContext] Extracted:', { access_token: access_token ? 'present' : 'missing', adminData })
        
        const validatedAdmin = validateAdminUser(adminData)
        console.log('[AuthContext] Validated admin:', validatedAdmin)
        
        if (validatedAdmin && access_token) {
          localStorage.setItem('adminToken', access_token)
          localStorage.setItem('adminData', JSON.stringify(validatedAdmin))
          console.log('[AuthContext] Stored in localStorage')
          
          setAdmin(validatedAdmin)
          setIsAuthenticated(true)
          setIsAdmin(true)
          setUser(null) // Clear user session when admin logs in
          console.log('[AuthContext] State updated - login successful!')
        } else {
          console.error('[AuthContext] Validation failed:', { validatedAdmin, access_token })
          throw new Error('Invalid admin data received')
        }
      } else {
        console.error('[AuthContext] Response not successful:', response)
        throw new Error(response.error || 'Admin login failed')
      }
    } catch (error) {
      console.error('[AuthContext] Login error:', error)
      const errorMessage = error instanceof Error ? error.message : 'Admin login failed'
      setError(errorMessage)
      throw error
    } finally {
      setIsLoading(false)
      console.log('[AuthContext] Login process complete')
    }
  }

  const adminLogout = (): void => {
    try {
      localStorage.removeItem('adminToken')
      localStorage.removeItem('adminData')
      
      setAdmin(null)
      setIsAuthenticated(false)
      setIsAdmin(false)
      setError(null)
    } catch (error) {
      console.error('Admin logout error:', error)
    }
  }

  // FIXED: Permission checking with dot-to-underscore conversion
  const hasPermission = (permission: string): boolean => {
    if (!isAdmin || !admin) return false
    
    try {
      // Ensure permissions is always an array
      const permissions = isValidArray(admin.permissions) ? admin.permissions : []
      
      // Check for wildcard permission
      if (permissions.includes('*' as AdminPermission)) return true
      
      // Convert underscore format to dot format for comparison
      const dotFormatPermission = permission.replace(/_/g, '.')
      
      // Check both formats: exact match and converted format
      return permissions.includes(permission as AdminPermission) || 
             permissions.includes(dotFormatPermission as AdminPermission)
    } catch (error) {
      console.error('Permission check error:', error)
      return false
    }
  }

  const hasRole = (role: string): boolean => {
    if (!isAdmin || !admin) return false
    
    try {
      return admin.role === role
    } catch (error) {
      console.error('Role check error:', error)
      return false
    }
  }

  // FIXED: Simplified canAccess function
  const canAccess = (requiredPermissions: string[] = [], requiredRoles: string[] = []): boolean => {
    if (!isAdmin || !admin) return false
    
    try {
      // Ensure parameters are arrays
      const permissions = isValidArray(requiredPermissions) ? requiredPermissions : []
      const roles = isValidArray(requiredRoles) ? requiredRoles : []
      
      // Check permissions - user must have ALL required permissions
      if (permissions.length > 0) {
        const hasAllPermissions = permissions.every(permission => hasPermission(permission))
        if (!hasAllPermissions) return false
      }
      
      // Check roles - user must have at least one of the required roles
      if (roles.length > 0) {
        const hasRequiredRole = roles.some(role => hasRole(role))
        if (!hasRequiredRole) return false
      }
      
      return true
    } catch (error) {
      console.error('Access check error:', error)
      return false
    }
  }

  const clearError = (): void => {
    setError(null)
  }

  const contextValue: AuthContextType = {
    user,
    admin,
    isAuthenticated,
    isAdmin,
    isLoading,
    error,
    login,
    register,
    logout,
    adminLogin,
    adminLogout,
    hasPermission,
    hasRole,
    canAccess,
    clearError
  }

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  )
}

export default AuthProvider

