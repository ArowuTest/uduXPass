import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { 
  Users,
  Search,
  Filter,
  Plus,
  Eye,
  Edit,
  Trash2,
  Download,
  RefreshCw,
  User,
  Mail,
  Phone,
  Calendar,
  Shield,
  UserCheck,
  UserX,
  Crown,
  MoreHorizontal
} from 'lucide-react'

// TypeScript interfaces
interface AdminUser {
  id: string
  firstName: string
  lastName: string
  email: string
  phone?: string
  role: 'admin' | 'super_admin' | 'moderator'
  status: 'active' | 'inactive' | 'suspended'
  createdAt: string
  lastLogin?: string
  permissions?: string[]
}

interface CreateUserData {
  firstName: string
  lastName: string
  email: string
  phone: string
  role: string
  password: string
  status: string
}

interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}

const AdminUserManagementPage: React.FC = () => {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [filteredUsers, setFilteredUsers] = useState<AdminUser[]>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [searchTerm, setSearchTerm] = useState<string>('')
  const [roleFilter, setRoleFilter] = useState<string>('all')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [selectedUser, setSelectedUser] = useState<AdminUser | null>(null)
  const [showUserDialog, setShowUserDialog] = useState<boolean>(false)
  const [showCreateUserDialog, setShowCreateUserDialog] = useState<boolean>(false)
  const [showEditUserDialog, setShowEditUserDialog] = useState<boolean>(false)
  const [editingUserId, setEditingUserId] = useState<string | null>(null)
  const [isExporting, setIsExporting] = useState<boolean>(false)
  const [newUserData, setNewUserData] = useState<CreateUserData>({
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    role: 'user',
    password: '',
    status: 'active'
  })

  useEffect(() => {
    fetchUsers()
  }, [])

  useEffect(() => {
    filterUsers()
  }, [users, searchTerm, roleFilter, statusFilter])

  const fetchUsers = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('/v1/admin/users', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: PaginatedResponse<AdminUser> } = await response.json()
        setUsers(result.data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch users:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const filterUsers = (): void => {
    let filtered = users

    if (searchTerm) {
      filtered = filtered.filter(user => 
        user.firstName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.lastName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.email.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (roleFilter !== 'all') {
      filtered = filtered.filter(user => user.role === roleFilter)
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter(user => user.status === statusFilter)
    }

    setFilteredUsers(filtered)
  }

  const handleViewUser = async (userId: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`/v1/admin/users/${userId}`, {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: AdminUser } = await response.json()
        setSelectedUser(result.data)
        setShowUserDialog(true)
      }
    } catch (error) {
      console.error('Failed to fetch user details:', error)
    }
  }

  const handleDeleteUser = async (userId: string): Promise<void> => {
    if (confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
      try {
        const adminToken = localStorage.getItem('adminToken')
        const response = await fetch(`/v1/admin/users/${userId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${adminToken}`,
            'Content-Type': 'application/json'
          }
        })
        
        if (response.ok) {
          fetchUsers()
          alert('User deleted successfully!')
        } else {
          alert('Failed to delete user')
        }
      } catch (error) {
        console.error('Failed to delete user:', error)
        alert('An error occurred while deleting the user.')
      }
    }
  }

  const handleUpdateUserStatus = async (userId: string, newStatus: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`/v1/admin/users/${userId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ status: newStatus })
      })
      
      if (response.ok) {
        fetchUsers()
        alert('User status updated successfully!')
      } else {
        alert('Failed to update user status')
      }
    } catch (error) {
      console.error('Failed to update user status:', error)
      alert('An error occurred while updating user status.')
    }
  }

  // NEW: Add User functionality
  const handleAddUser = (): void => {
    setNewUserData({
      firstName: '',
      lastName: '',
      email: '',
      phone: '',
      role: 'user',
      password: '',
      status: 'active'
    })
    setShowCreateUserDialog(true)
  }

  const handleCreateUser = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault()
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('/v1/admin/users', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newUserData)
      })
      
      if (response.ok) {
        fetchUsers()
        setShowCreateUserDialog(false)
        alert('User created successfully!')
      } else {
        const error = await response.text()
        alert('Failed to create user: ' + error)
      }
    } catch (error) {
      console.error('Failed to create user:', error)
      alert('An error occurred while creating the user.')
    }
  }

  // NEW: Edit User functionality
  const handleEditUser = async (userId: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`/v1/admin/users/${userId}`, {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: AdminUser } = await response.json()
        setNewUserData({
          firstName: result.data.firstName,
          lastName: result.data.lastName,
          email: result.data.email,
          phone: result.data.phone || '',
          role: result.data.role,
          password: '', // Don't pre-fill password
          status: result.data.status
        })
        setEditingUserId(userId)
        setShowEditUserDialog(true)
      }
    } catch (error) {
      console.error('Failed to fetch user for editing:', error)
    }
  }

  const handleUpdateUser = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault()
    try {
      const adminToken = localStorage.getItem('adminToken')
      const updateData = { ...newUserData }
      
      // Remove password if it's empty (don't update password)
      if (!updateData.password) {
        delete (updateData as any).password
      }
      
      const response = await fetch(`/v1/admin/users/${editingUserId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateData)
      })
      
      if (response.ok) {
        fetchUsers()
        setShowEditUserDialog(false)
        setEditingUserId(null)
        alert('User updated successfully!')
      } else {
        const error = await response.text()
        alert('Failed to update user: ' + error)
      }
    } catch (error) {
      console.error('Failed to update user:', error)
      alert('An error occurred while updating the user.')
    }
  }

  // NEW: Export Users functionality
  const handleExportUsers = async (): Promise<void> => {
    setIsExporting(true)
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('/v1/admin/users/export', {
        headers: {
          'Authorization': `Bearer ${adminToken}`
        }
      })
      
      if (response.ok) {
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `users-export-${new Date().toISOString().split('T')[0]}.csv`
        a.click()
        window.URL.revokeObjectURL(url)
        alert('Users exported successfully!')
      } else {
        alert('Failed to export users')
      }
    } catch (error) {
      console.error('Export failed:', error)
      alert('An error occurred while exporting users.')
    } finally {
      setIsExporting(false)
    }
  }

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString('en-NG', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getRoleBadge = (role: string): JSX.Element => {
    const roleConfig: Record<string, { color: string; label: string; icon: React.ComponentType<any> }> = {
      super_admin: { color: 'bg-purple-100 text-purple-800', label: 'Super Admin', icon: Crown },
      admin: { color: 'bg-blue-100 text-blue-800', label: 'Admin', icon: Shield },
      moderator: { color: 'bg-green-100 text-green-800', label: 'Moderator', icon: UserCheck }
    }
    
    const config = roleConfig[role] || { color: 'bg-gray-100 text-gray-800', label: role, icon: User }
    const Icon = config.icon
    
    return (
      <Badge className={`${config.color} flex items-center gap-1`}>
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getStatusBadge = (status: string): JSX.Element => {
    const statusConfig: Record<string, { color: string; label: string; icon: React.ComponentType<any> }> = {
      active: { color: 'bg-green-100 text-green-800', label: 'Active', icon: UserCheck },
      inactive: { color: 'bg-gray-100 text-gray-800', label: 'Inactive', icon: UserX },
      suspended: { color: 'bg-red-100 text-red-800', label: 'Suspended', icon: UserX }
    }
    
    const config = statusConfig[status] || { color: 'bg-gray-100 text-gray-800', label: status, icon: User }
    const Icon = config.icon
    
    return (
      <Badge className={`${config.color} flex items-center gap-1`}>
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getUserStats = () => {
    const total = users.length
    const active = users.filter(u => u.status === 'active').length
    const admins = users.filter(u => u.role === 'admin').length
    const superAdmins = users.filter(u => u.role === 'super_admin').length
    const suspended = users.filter(u => u.status === 'suspended').length
    
    return { total, active, admins, superAdmins, suspended }
  }

  const stats = getUserStats()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-2">Loading users...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">User Management</h1>
          <p className="text-muted-foreground">
            Manage regular users and customers
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" onClick={fetchUsers}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button 
            variant="outline" 
            onClick={handleExportUsers}
            disabled={isExporting}
          >
            <Download className="h-4 w-4 mr-2" />
            {isExporting ? 'Exporting...' : 'Export Users'}
          </Button>
          <Button onClick={handleAddUser}>
            <Plus className="h-4 w-4 mr-2" />
            Add User
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Users</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-xs text-muted-foreground">
              All admin users
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active</CardTitle>
            <UserCheck className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{stats.active}</div>
            <p className="text-xs text-muted-foreground">
              Currently active
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Admins</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">{stats.admins}</div>
            <p className="text-xs text-muted-foreground">
              Admin role
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Super Admins</CardTitle>
            <Crown className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">{stats.superAdmins}</div>
            <p className="text-xs text-muted-foreground">
              Super admin role
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Suspended</CardTitle>
            <UserX className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{stats.suspended}</div>
            <p className="text-xs text-muted-foreground">
              Suspended users
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center space-x-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search users by name or email..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <Select value={roleFilter} onValueChange={setRoleFilter}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Filter by role" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Roles</SelectItem>
                <SelectItem value="super_admin">Super Admin</SelectItem>
                <SelectItem value="admin">Admin</SelectItem>
                <SelectItem value="moderator">Moderator</SelectItem>
              </SelectContent>
            </Select>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="inactive">Inactive</SelectItem>
                <SelectItem value="suspended">Suspended</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Users Table */}
      <Card>
        <CardHeader>
          <CardTitle>Users ({filteredUsers.length})</CardTitle>
          <CardDescription>
            Manage admin users and their access permissions
          </CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">User</th>
                  <th className="text-left p-4">Contact</th>
                  <th className="text-left p-4">Role</th>
                  <th className="text-left p-4">Status</th>
                  <th className="text-left p-4">Created</th>
                  <th className="text-left p-4">Last Login</th>
                  <th className="text-left p-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredUsers.map((user) => (
                  <tr key={user.id} className="border-b hover:bg-gray-50">
                    <td className="p-4">
                      <div className="flex items-center space-x-3">
                        <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                          <User className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <div className="font-medium">{user.firstName} {user.lastName}</div>
                          <div className="text-sm text-muted-foreground">ID: {user.id}</div>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="space-y-1">
                        <div className="text-sm flex items-center">
                          <Mail className="h-3 w-3 mr-1 text-muted-foreground" />
                          {user.email}
                        </div>
                        {user.phone && (
                          <div className="text-sm flex items-center">
                            <Phone className="h-3 w-3 mr-1 text-muted-foreground" />
                            {user.phone}
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="p-4">{getRoleBadge(user.role)}</td>
                    <td className="p-4">{getStatusBadge(user.status)}</td>
                    <td className="p-4">
                      <div className="flex items-center text-sm">
                        <Calendar className="h-3 w-3 mr-1 text-muted-foreground" />
                        {formatDate(user.createdAt)}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="text-sm">
                        {user.lastLogin ? formatDate(user.lastLogin) : 'Never'}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center space-x-2">
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => handleViewUser(user.id)}
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => handleEditUser(user.id)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Select onValueChange={(value) => {
                          if (value === 'delete') {
                            handleDeleteUser(user.id)
                          } else {
                            handleUpdateUserStatus(user.id, value)
                          }
                        }}>
                          <SelectTrigger className="w-8 h-8 p-0">
                            <MoreHorizontal className="h-4 w-4" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="active">Set Active</SelectItem>
                            <SelectItem value="inactive">Set Inactive</SelectItem>
                            <SelectItem value="suspended">Suspend User</SelectItem>
                            <SelectItem value="delete">Delete User</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Create User Dialog */}
      <Dialog open={showCreateUserDialog} onOpenChange={setShowCreateUserDialog}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Add New User</DialogTitle>
            <DialogDescription>
              Create a new admin user account
            </DialogDescription>
          </DialogHeader>
          
          <form onSubmit={handleCreateUser} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="firstName">First Name</Label>
                <Input
                  id="firstName"
                  required
                  value={newUserData.firstName}
                  onChange={(e) => setNewUserData({...newUserData, firstName: e.target.value})}
                />
              </div>
              <div>
                <Label htmlFor="lastName">Last Name</Label>
                <Input
                  id="lastName"
                  required
                  value={newUserData.lastName}
                  onChange={(e) => setNewUserData({...newUserData, lastName: e.target.value})}
                />
              </div>
            </div>
            
            <div>
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                required
                value={newUserData.email}
                onChange={(e) => setNewUserData({...newUserData, email: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="phone">Phone Number</Label>
              <Input
                id="phone"
                value={newUserData.phone}
                onChange={(e) => setNewUserData({...newUserData, phone: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                required
                value={newUserData.password}
                onChange={(e) => setNewUserData({...newUserData, password: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="role">Role</Label>
              <Select value={newUserData.role} onValueChange={(value) => setNewUserData({...newUserData, role: value})}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">User</SelectItem>
                  <SelectItem value="moderator">Moderator</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div>
              <Label htmlFor="status">Status</Label>
              <Select value={newUserData.status} onValueChange={(value) => setNewUserData({...newUserData, status: value})}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="flex justify-end space-x-2">
              <Button type="button" variant="outline" onClick={() => setShowCreateUserDialog(false)}>
                Cancel
              </Button>
              <Button type="submit">Create User</Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Edit User Dialog */}
      <Dialog open={showEditUserDialog} onOpenChange={setShowEditUserDialog}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Edit User</DialogTitle>
            <DialogDescription>
              Update user information
            </DialogDescription>
          </DialogHeader>
          
          <form onSubmit={handleUpdateUser} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="editFirstName">First Name</Label>
                <Input
                  id="editFirstName"
                  required
                  value={newUserData.firstName}
                  onChange={(e) => setNewUserData({...newUserData, firstName: e.target.value})}
                />
              </div>
              <div>
                <Label htmlFor="editLastName">Last Name</Label>
                <Input
                  id="editLastName"
                  required
                  value={newUserData.lastName}
                  onChange={(e) => setNewUserData({...newUserData, lastName: e.target.value})}
                />
              </div>
            </div>
            
            <div>
              <Label htmlFor="editEmail">Email</Label>
              <Input
                id="editEmail"
                type="email"
                required
                value={newUserData.email}
                onChange={(e) => setNewUserData({...newUserData, email: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editPhoneNumber">Phone Number</Label>
              <Input
                id="editPhoneNumber"
                value={newUserData.phone}
                onChange={(e) => setNewUserData({...newUserData, phone: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editPassword">Password (leave blank to keep current)</Label>
              <Input
                id="editPassword"
                type="password"
                value={newUserData.password}
                onChange={(e) => setNewUserData({...newUserData, password: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editRole">Role</Label>
              <Select value={newUserData.role} onValueChange={(value) => setNewUserData({...newUserData, role: value})}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">User</SelectItem>
                  <SelectItem value="moderator">Moderator</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div>
              <Label htmlFor="editStatus">Status</Label>
              <Select value={newUserData.status} onValueChange={(value) => setNewUserData({...newUserData, status: value})}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                  <SelectItem value="suspended">Suspended</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="flex justify-end space-x-2">
              <Button type="button" variant="outline" onClick={() => setShowEditUserDialog(false)}>
                Cancel
              </Button>
              <Button type="submit">Update User</Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* User Details Dialog */}
      <Dialog open={showUserDialog} onOpenChange={setShowUserDialog}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>User Details</DialogTitle>
            <DialogDescription>
              Complete information about the selected user
            </DialogDescription>
          </DialogHeader>
          
          {selectedUser && (
            <div className="space-y-6">
              <div className="flex items-center space-x-4">
                <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
                  <User className="h-8 w-8 text-blue-600" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold">{selectedUser.firstName} {selectedUser.lastName}</h3>
                  <p className="text-muted-foreground">{selectedUser.email}</p>
                  <div className="flex items-center space-x-2 mt-1">
                    {getRoleBadge(selectedUser.role)}
                    {getStatusBadge(selectedUser.status)}
                  </div>
                </div>
              </div>

              <Tabs defaultValue="details" className="w-full">
                <TabsList className="grid w-full grid-cols-3">
                  <TabsTrigger value="details">Details</TabsTrigger>
                  <TabsTrigger value="permissions">Permissions</TabsTrigger>
                  <TabsTrigger value="activity">Activity</TabsTrigger>
                </TabsList>
                
                <TabsContent value="details" className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <Label>User ID</Label>
                      <p className="text-sm">{selectedUser.id}</p>
                    </div>
                    <div>
                      <Label>Phone Number</Label>
                      <p className="text-sm">{selectedUser.phone || 'Not provided'}</p>
                    </div>
                    <div>
                      <Label>Created</Label>
                      <p className="text-sm">{formatDate(selectedUser.createdAt)}</p>
                    </div>
                    <div>
                      <Label>Last Login</Label>
                      <p className="text-sm">
                        {selectedUser.lastLogin ? formatDate(selectedUser.lastLogin) : 'Never'}
                      </p>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="permissions">
                  <div className="space-y-2">
                    {selectedUser.permissions && selectedUser.permissions.length > 0 ? (
                      selectedUser.permissions.map((permission, index) => (
                        <Badge key={index} variant="outline">{permission}</Badge>
                      ))
                    ) : (
                      <p className="text-sm text-muted-foreground">No specific permissions assigned</p>
                    )}
                  </div>
                </TabsContent>
                
                <TabsContent value="activity">
                  <p className="text-sm text-muted-foreground">User activity log will be displayed here.</p>
                </TabsContent>
              </Tabs>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default AdminUserManagementPage

