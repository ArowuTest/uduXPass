import React, { useState, useEffect } from 'react'
import { scannersAPI } from '@/services/api';
import { toast } from "@/components/ui/toaster";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { 
  Scan,
  Search,
  Filter,
  Plus,
  Eye,
  Edit,
  Trash2,
  Download,
  RefreshCw,
  Smartphone,
  Wifi,
  WifiOff,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Calendar,
  User,
  MapPin,
  Activity,
  MoreHorizontal,
  Loader2
} from 'lucide-react'

// TypeScript interfaces
interface Scanner {
  id: string
  name: string
  description?: string
  location: string
  assignedUser?: string
  eventId?: string
  eventName?: string
  status: 'active' | 'inactive' | 'maintenance' | 'offline'
  isOnline: boolean
  totalScans?: number
  scansToday?: number
  lastActivity?: string
  createdAt: string
}

interface CreateScannerData {
  name: string
  description: string
  location: string
  assignedUser: string
  eventId: string
  status: string
}

const AdminScannerManagementPage: React.FC = () => {
  
  const [scanners, setScanners] = useState<Scanner[]>([])
  const [filteredScanners, setFilteredScanners] = useState<Scanner[]>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [isCreating, setIsCreating] = useState<boolean>(false)
  const [isUpdating, setIsUpdating] = useState<boolean>(false)
  const [searchTerm, setSearchTerm] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [selectedScanner, setSelectedScanner] = useState<Scanner | null>(null)
  const [showScannerDialog, setShowScannerDialog] = useState<boolean>(false)
  const [showCreateScannerDialog, setShowCreateScannerDialog] = useState<boolean>(false)
  const [showEditScannerDialog, setShowEditScannerDialog] = useState<boolean>(false)
  const [editingScannerId, setEditingScannerId] = useState<string | null>(null)
  const [newScannerData, setNewScannerData] = useState<CreateScannerData>({
    name: '',
    description: '',
    location: '',
    assignedUser: '',
    eventId: '',
    status: 'active'
  })

  useEffect(() => {
    fetchScanners()
  }, [])

  useEffect(() => {
    filterScanners()
  }, [scanners, searchTerm, statusFilter])

  const fetchScanners = async (): Promise<void> => {
    try {
      setIsLoading(true)
      const response = await scannersAPI.getAll()
      
      if (response.success) {
        setScanners(response.data)
      } else {
        toast({
          title: "Error",
          description: response.error || "Failed to fetch scanners",
          variant: "destructive"
        })
      }
    } catch (error) {
      console.error('Failed to fetch scanners:', error)
      toast({
        title: "Error",
        description: "An error occurred while fetching scanners",
        variant: "destructive"
      })
    } finally {
      setIsLoading(false)
    }
  }

  const filterScanners = (): void => {
    if (!Array.isArray(scanners)) {
      setFilteredScanners([])
      return
    }

    let filtered = scanners

    if (searchTerm) {
      filtered = filtered.filter(scanner =>
        scanner.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        scanner.location.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (scanner.assignedUser && scanner.assignedUser.toLowerCase().includes(searchTerm.toLowerCase()))
      )
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter(scanner => scanner.status === statusFilter)
    }

    setFilteredScanners(filtered)
  }

  const handleAddScanner = (): void => {
    setNewScannerData({
      name: '',
      description: '',
      location: '',
      assignedUser: '',
      eventId: '',
      status: 'active'
    })
    setShowCreateScannerDialog(true)
  }

  const handleCreateScanner = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault()
    
    // Form validation
    if (!newScannerData.name.trim()) {
      toast({
        title: "Validation Error",
        description: "Scanner name is required",
        variant: "destructive"
      })
      return
    }
    
    if (!newScannerData.location.trim()) {
      toast({
        title: "Validation Error", 
        description: "Location is required",
        variant: "destructive"
      })
      return
    }

    try {
      setIsCreating(true)
      const response = await scannersAPI.create(newScannerData)
      
      if (response.success) {
        await fetchScanners()
        setShowCreateScannerDialog(false)
        toast({
          title: "Success",
          description: "Scanner created successfully!"
        })
      } else {
        toast({
          title: "Error",
          description: response.error || "Failed to create scanner",
          variant: "destructive"
        })
      }
    } catch (error) {
      console.error('Failed to create scanner:', error)
      toast({
        title: "Error",
        description: "An error occurred while creating the scanner",
        variant: "destructive"
      })
    } finally {
      setIsCreating(false)
    }
  }

  const handleEditScanner = async (scannerId: string): Promise<void> => {
    try {
      const response = await scannersAPI.getById(scannerId)
      
      if (response.success) {
        setNewScannerData({
          name: response.data.name,
          description: response.data.description || '',
          location: response.data.location,
          assignedUser: response.data.assignedUser || '',
          eventId: response.data.eventId || '',
          status: response.data.status
        })
        setEditingScannerId(scannerId)
        setShowEditScannerDialog(true)
      } else {
        toast({
          title: "Error",
          description: response.error || "Failed to fetch scanner details",
          variant: "destructive"
        })
      }
    } catch (error) {
      console.error('Failed to fetch scanner for editing:', error)
      toast({
        title: "Error",
        description: "An error occurred while fetching scanner details",
        variant: "destructive"
      })
    }
  }

  const handleUpdateScanner = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault()
    
    if (!editingScannerId) return

    // Form validation
    if (!newScannerData.name.trim()) {
      toast({
        title: "Validation Error",
        description: "Scanner name is required",
        variant: "destructive"
      })
      return
    }

    try {
      setIsUpdating(true)
      const response = await scannersAPI.update(editingScannerId, newScannerData)
      
      if (response.success) {
        await fetchScanners()
        setShowEditScannerDialog(false)
        setEditingScannerId(null)
        toast({
          title: "Success",
          description: "Scanner updated successfully!"
        })
      } else {
        toast({
          title: "Error",
          description: response.error || "Failed to update scanner",
          variant: "destructive"
        })
      }
    } catch (error) {
      console.error('Failed to update scanner:', error)
      toast({
        title: "Error",
        description: "An error occurred while updating the scanner",
        variant: "destructive"
      })
    } finally {
      setIsUpdating(false)
    }
  }

  const handleDeleteScanner = async (scannerId: string): Promise<void> => {
    if (!confirm('Are you sure you want to delete this scanner?')) return

    try {
      const response = await scannersAPI.delete(scannerId)
      
      if (response.success) {
        await fetchScanners()
        toast({
          title: "Success",
          description: "Scanner deleted successfully!"
        })
      } else {
        toast({
          title: "Error",
          description: response.error || "Failed to delete scanner",
          variant: "destructive"
        })
      }
    } catch (error) {
      console.error('Failed to delete scanner:', error)
      toast({
        title: "Error",
        description: "An error occurred while deleting the scanner",
        variant: "destructive"
      })
    }
  }

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      active: { color: 'bg-green-100 text-green-800', label: 'Active' },
      inactive: { color: 'bg-gray-100 text-gray-800', label: 'Inactive' },
      maintenance: { color: 'bg-yellow-100 text-yellow-800', label: 'Maintenance' },
      offline: { color: 'bg-red-100 text-red-800', label: 'Offline' }
    }
    
    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.inactive
    return <Badge className={config.color}>{config.label}</Badge>
  }

  const getConnectionBadge = (isOnline: boolean) => {
    return isOnline ? (
      <Badge className="bg-green-100 text-green-800">
        <Wifi className="w-3 h-3 mr-1" />
        Online
      </Badge>
    ) : (
      <Badge className="bg-red-100 text-red-800">
        <WifiOff className="w-3 h-3 mr-1" />
        Offline
      </Badge>
    )
  }

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString()
  }

  const stats = {
    total: Array.isArray(scanners) ? scanners.length : 0,
    active: Array.isArray(scanners) ? scanners.filter(s => s.status === 'active').length : 0,
    offline: Array.isArray(scanners) ? scanners.filter(s => !s.isOnline).length : 0,
    maintenance: Array.isArray(scanners) ? scanners.filter(s => s.status === 'maintenance').length : 0,
    totalScans: Array.isArray(scanners) ? scanners.reduce((sum, s) => sum + (s.totalScans || 0), 0) : 0
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin" />
        <span className="ml-2">Loading scanners...</span>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Scanner Management</h1>
          <p className="text-muted-foreground">
            Manage scanning devices and monitor their performance
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" onClick={fetchScanners} disabled={isLoading}>
            <RefreshCw className="w-4 h-4 mr-2" />
            Refresh
          </Button>
          <Button onClick={handleAddScanner}>
            <Plus className="w-4 h-4 mr-2" />
            Add Scanner
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total</CardTitle>
            <Scan className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-xs text-muted-foreground">
              Total scanners
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active</CardTitle>
            <CheckCircle className="h-4 w-4 text-muted-foreground" />
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
            <CardTitle className="text-sm font-medium">Offline</CardTitle>
            <WifiOff className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{stats.offline}</div>
            <p className="text-xs text-muted-foreground">
              Not responding
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Maintenance</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{stats.maintenance}</div>
            <p className="text-xs text-muted-foreground">
              Under maintenance
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Scans</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">{stats.totalScans}</div>
            <p className="text-xs text-muted-foreground">
              All time scans
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
                  placeholder="Search scanners by name, location, or assigned user..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="inactive">Inactive</SelectItem>
                <SelectItem value="maintenance">Maintenance</SelectItem>
                <SelectItem value="offline">Offline</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Scanners Table */}
      <Card>
        <CardHeader>
          <CardTitle>Scanners ({filteredScanners.length})</CardTitle>
          <CardDescription>
            Manage scanning devices and monitor their performance
          </CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">Scanner</th>
                  <th className="text-left p-4">Location</th>
                  <th className="text-left p-4">Assigned User</th>
                  <th className="text-left p-4">Status</th>
                  <th className="text-left p-4">Connection</th>
                  <th className="text-left p-4">Scans Today</th>
                  <th className="text-left p-4">Last Activity</th>
                  <th className="text-left p-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredScanners.map((scanner) => (
                  <tr key={scanner.id} className="border-b hover:bg-muted/50">
                    <td className="p-4">
                      <div className="flex items-center space-x-3">
                        <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                          <Smartphone className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <div className="font-medium">{scanner.name}</div>
                          <div className="text-sm text-muted-foreground">{scanner.description}</div>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center">
                        <MapPin className="h-4 w-4 mr-1 text-muted-foreground" />
                        {scanner.location}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center">
                        <User className="h-4 w-4 mr-1 text-muted-foreground" />
                        {scanner.assignedUser || 'Unassigned'}
                      </div>
                    </td>
                    <td className="p-4">
                      {getStatusBadge(scanner.status)}
                    </td>
                    <td className="p-4">
                      {getConnectionBadge(scanner.isOnline)}
                    </td>
                    <td className="p-4">
                      <div className="font-medium">{scanner.scansToday || 0}</div>
                    </td>
                    <td className="p-4">
                      <div className="text-sm">
                        {scanner.lastActivity ? formatDate(scanner.lastActivity) : 'Never'}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center space-x-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => {
                            setSelectedScanner(scanner)
                            setShowScannerDialog(true)
                          }}
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleEditScanner(scanner.id)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDeleteScanner(scanner.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Create Scanner Dialog */}
      <Dialog open={showCreateScannerDialog} onOpenChange={setShowCreateScannerDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add New Scanner</DialogTitle>
            <DialogDescription>
              Create a new scanner device for event ticket validation
            </DialogDescription>
          </DialogHeader>
          
          <form onSubmit={handleCreateScanner} className="space-y-4">
            <div>
              <Label htmlFor="name">Scanner Name *</Label>
              <Input
                id="name"
                required
                value={newScannerData.name}
                onChange={(e) => setNewScannerData({...newScannerData, name: e.target.value})}
                placeholder="e.g., Main Entrance Scanner"
              />
            </div>
            
            <div>
              <Label htmlFor="description">Description</Label>
              <Input
                id="description"
                value={newScannerData.description}
                onChange={(e) => setNewScannerData({...newScannerData, description: e.target.value})}
                placeholder="Optional description"
              />
            </div>
            
            <div>
              <Label htmlFor="location">Location *</Label>
              <Input
                id="location"
                required
                value={newScannerData.location}
                onChange={(e) => setNewScannerData({...newScannerData, location: e.target.value})}
                placeholder="e.g., Gate A, VIP Entrance"
              />
            </div>
            
            <div>
              <Label htmlFor="assignedUser">Assigned User</Label>
              <Input
                id="assignedUser"
                value={newScannerData.assignedUser}
                onChange={(e) => setNewScannerData({...newScannerData, assignedUser: e.target.value})}
                placeholder="Scanner operator username"
              />
            </div>
            
            <div>
              <Label htmlFor="eventId">Event (Optional)</Label>
              <Input
                id="eventId"
                value={newScannerData.eventId}
                onChange={(e) => setNewScannerData({...newScannerData, eventId: e.target.value})}
                placeholder="Event ID to assign scanner to"
              />
            </div>
            
            <div>
              <Label htmlFor="status">Status</Label>
              <Select value={newScannerData.status} onValueChange={(value) => setNewScannerData({...newScannerData, status: value})}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                  <SelectItem value="maintenance">Maintenance</SelectItem>
                  <SelectItem value="offline">Offline</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="flex justify-end space-x-2">
              <Button type="button" variant="outline" onClick={() => setShowCreateScannerDialog(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isCreating}>
                {isCreating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Create Scanner
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Edit Scanner Dialog */}
      <Dialog open={showEditScannerDialog} onOpenChange={setShowEditScannerDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Scanner</DialogTitle>
            <DialogDescription>
              Update scanner device information
            </DialogDescription>
          </DialogHeader>
          
          <form onSubmit={handleUpdateScanner} className="space-y-4">
            <div>
              <Label htmlFor="editName">Scanner Name *</Label>
              <Input
                id="editName"
                required
                value={newScannerData.name}
                onChange={(e) => setNewScannerData({...newScannerData, name: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editDescription">Description</Label>
              <Input
                id="editDescription"
                value={newScannerData.description}
                onChange={(e) => setNewScannerData({...newScannerData, description: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editLocation">Location *</Label>
              <Input
                id="editLocation"
                required
                value={newScannerData.location}
                onChange={(e) => setNewScannerData({...newScannerData, location: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editAssignedUser">Assigned User</Label>
              <Input
                id="editAssignedUser"
                value={newScannerData.assignedUser}
                onChange={(e) => setNewScannerData({...newScannerData, assignedUser: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editEventId">Event (Optional)</Label>
              <Input
                id="editEventId"
                value={newScannerData.eventId}
                onChange={(e) => setNewScannerData({...newScannerData, eventId: e.target.value})}
              />
            </div>
            
            <div>
              <Label htmlFor="editStatus">Status</Label>
              <Select value={newScannerData.status} onValueChange={(value) => setNewScannerData({...newScannerData, status: value})}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                  <SelectItem value="maintenance">Maintenance</SelectItem>
                  <SelectItem value="offline">Offline</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="flex justify-end space-x-2">
              <Button type="button" variant="outline" onClick={() => setShowEditScannerDialog(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isUpdating}>
                {isUpdating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Update Scanner
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Scanner Details Dialog */}
      <Dialog open={showScannerDialog} onOpenChange={setShowScannerDialog}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Scanner Details</DialogTitle>
            <DialogDescription>
              Complete information about the selected scanner
            </DialogDescription>
          </DialogHeader>
          
          {selectedScanner && (
            <div className="space-y-6">
              <div className="flex items-center space-x-4">
                <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
                  <Smartphone className="h-8 w-8 text-blue-600" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold">{selectedScanner.name}</h3>
                  <p className="text-muted-foreground">{selectedScanner.location}</p>
                  <div className="flex items-center space-x-2 mt-1">
                    {getStatusBadge(selectedScanner.status)}
                    {getConnectionBadge(selectedScanner.isOnline)}
                  </div>
                </div>
              </div>

              <Tabs defaultValue="details" className="w-full">
                <TabsList className="grid w-full grid-cols-3">
                  <TabsTrigger value="details">Details</TabsTrigger>
                  <TabsTrigger value="activity">Activity</TabsTrigger>
                  <TabsTrigger value="settings">Settings</TabsTrigger>
                </TabsList>
                
                <TabsContent value="details" className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <Label>Scanner ID</Label>
                      <p className="text-sm">{selectedScanner.id}</p>
                    </div>
                    <div>
                      <Label>Description</Label>
                      <p className="text-sm">{selectedScanner.description || 'No description'}</p>
                    </div>
                    <div>
                      <Label>Assigned User</Label>
                      <p className="text-sm">{selectedScanner.assignedUser || 'Unassigned'}</p>
                    </div>
                    <div>
                      <Label>Event</Label>
                      <p className="text-sm">{selectedScanner.eventName || 'Not assigned to event'}</p>
                    </div>
                    <div>
                      <Label>Total Scans</Label>
                      <p className="text-sm">{selectedScanner.totalScans || 0}</p>
                    </div>
                    <div>
                      <Label>Last Activity</Label>
                      <p className="text-sm">
                        {selectedScanner.lastActivity ? formatDate(selectedScanner.lastActivity) : 'Never'}
                      </p>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="activity">
                  <p className="text-sm text-muted-foreground">Scanner activity log will be displayed here.</p>
                </TabsContent>
                
                <TabsContent value="settings">
                  <p className="text-sm text-muted-foreground">Scanner configuration settings will be displayed here.</p>
                </TabsContent>
              </Tabs>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default AdminScannerManagementPage
