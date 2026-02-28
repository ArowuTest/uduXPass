import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { 
  Calendar, 
  MapPin, 
  Users, 
  DollarSign, 
  Plus,
  Search,
  Filter,
  Edit,
  Eye,
  Trash2,
  RefreshCw,
  MoreHorizontal
} from 'lucide-react'

// TypeScript interfaces
interface Event {
  id: string
  title?: string
  name?: string // Backend might use 'name' instead of 'title'
  organizer?: string
  date?: string
  startDate?: string
  venue: string | { name: string; city: string }
  status: string
  ticketsSold?: number
  revenue?: number
  capacity?: number
  description?: string
  category?: string
}

const AdminEventsPage: React.FC = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  const [error, setError] = useState<string | null>(null)

  // Mock data for demonstration - will be replaced by API data
  const mockEvents: Event[] = [
    {
      id: 'event-1',
      title: 'Davido Live in Lagos',
      name: 'Davido Live in Lagos',
      organizer: 'uduXPass',
      startDate: '2025-12-15T18:00:00Z',
      venue: { name: 'Tafawa Balewa Square', city: 'Lagos' },
      status: 'published',
      ticketsSold: 350,
      revenue: 5250000,
      capacity: 1000,
      description: "Davido's exclusive concert featuring special guests",
      category: 'Music'
    },
    {
      id: 'event-2', 
      title: 'Burna Boy Concert',
      name: 'Burna Boy Concert',
      organizer: 'uduXPass',
      startDate: '2025-11-20T19:00:00Z',
      venue: { name: 'Eko Convention Centre', city: 'Lagos' },
      status: 'published',
      ticketsSold: 320,
      revenue: 4800000,
      capacity: 800,
      description: 'African Giant live performance',
      category: 'Music'
    },
    {
      id: 'event-3',
      title: 'Wizkid Live in Lagos',
      name: 'Wizkid Live in Lagos', 
      organizer: 'uduXPass',
      startDate: '2025-10-18T20:00:00Z',
      venue: { name: 'Eko Atlantic City', city: 'Lagos' },
      status: 'upcoming',
      ticketsSold: 0,
      revenue: 0,
      capacity: 1700,
      description: 'The Lagos grand finale of the Wizkid tour featuring special guest appearances',
      category: 'Music'
    }
  ]

  useEffect(() => {
    fetchEvents()
  }, [])

  const fetchEvents = async () => {
    try {
      setLoading(true)
      setError(null)
      
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('/v1/admin/events', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: { events: Event[], pagination: any }, success: boolean } = await response.json()
        console.log('Events API response:', result)
        
        // If API returns events, use them; otherwise use mock data
        if (result.success && result.data && result.data.events) {
          setEvents(result.data.events)
        } else {
          console.log('Using mock events data')
          setEvents(mockEvents)
        }
      } else {
        console.log('API failed, using mock events data')
        setEvents(mockEvents)
      }
    } catch (error) {
      console.error('Error fetching events:', error)
      console.log('Using mock events data due to error')
      setEvents(mockEvents)
    } finally {
      setLoading(false)
    }
  }

  const filterEvents = () => {
    let filtered = events

    if (searchTerm) {
      filtered = filtered.filter(event => {
        const eventTitle = event.title || event.name || ''
        const eventOrganizer = event.organizer || ''
        
        let venueText = ''
        if (typeof event.venue === 'string') {
          venueText = event.venue
        } else if (event.venue && event.venue.name && event.venue.city) {
          venueText = `${event.venue.name} ${event.venue.city}`
        }
        
        return (
          eventTitle.toLowerCase().includes(searchTerm.toLowerCase()) ||
          eventOrganizer.toLowerCase().includes(searchTerm.toLowerCase()) ||
          venueText.toLowerCase().includes(searchTerm.toLowerCase())
        )
      })
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter(event => event.status === statusFilter)
    }

    return filtered
  }

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      published: { label: 'Published', variant: 'default' as const },
      draft: { label: 'Draft', variant: 'secondary' as const },
      upcoming: { label: 'Upcoming', variant: 'outline' as const },
      completed: { label: 'Completed', variant: 'secondary' as const },
      cancelled: { label: 'Cancelled', variant: 'destructive' as const }
    }
    
    const config = statusConfig[status as keyof typeof statusConfig] || { label: status, variant: 'outline' as const }
    return <Badge variant={config.variant}>{config.label}</Badge>
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: 'NGN',
      minimumFractionDigits: 0
    }).format(amount)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getVenueText = (venue: string | { name: string; city: string }) => {
    if (typeof venue === 'string') {
      return venue
    }
    return `${venue.name}, ${venue.city}`
  }

  const filteredEvents = filterEvents()

  // Calculate summary statistics
  const totalEvents = events.length
  const publishedEvents = events.filter(e => e.status === 'published').length
  const draftEvents = events.filter(e => e.status === 'draft').length
  const totalRevenue = events.reduce((sum, event) => sum + (event.revenue || 0), 0)

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <RefreshCw className="h-8 w-8 animate-spin" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Events</h1>
          <p className="text-muted-foreground">Manage and monitor all events on the platform</p>
        </div>
        <Button onClick={() => navigate('/admin/events/create')}>
          <Plus className="h-4 w-4 mr-2" />
          Create Event
        </Button>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Events</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{totalEvents}</div>
            <p className="text-xs text-muted-foreground">All events on platform</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Published</CardTitle>
            <Eye className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{publishedEvents}</div>
            <p className="text-xs text-muted-foreground">Live events</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Draft</CardTitle>
            <Edit className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{draftEvents}</div>
            <p className="text-xs text-muted-foreground">Unpublished events</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalRevenue)}</div>
            <p className="text-xs text-muted-foreground">From all events</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex gap-4 items-center">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
          <Input
            placeholder="Search events by title, organizer, or venue..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10"
          />
        </div>
        
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger className="w-[180px]">
            <Filter className="h-4 w-4 mr-2" />
            <SelectValue placeholder="All Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="published">Published</SelectItem>
            <SelectItem value="draft">Draft</SelectItem>
            <SelectItem value="upcoming">Upcoming</SelectItem>
            <SelectItem value="completed">Completed</SelectItem>
            <SelectItem value="cancelled">Cancelled</SelectItem>
          </SelectContent>
        </Select>

        <Button variant="outline" onClick={fetchEvents}>
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      {/* Events Grid */}
      {error && (
        <div className="text-red-600 bg-red-50 p-4 rounded-md">
          Error: {error}
        </div>
      )}

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {filteredEvents.map((event) => (
          <Card key={event.id} className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <div className="flex justify-between items-start">
                <div className="space-y-1">
                  <CardTitle className="text-lg">{event.title || event.name}</CardTitle>
                  <CardDescription className="flex items-center">
                    <Calendar className="h-4 w-4 mr-1" />
                    {event.startDate ? formatDate(event.startDate) : 'Invalid Date'}
                  </CardDescription>
                </div>
                {getStatusBadge(event.status)}
              </div>
            </CardHeader>
            
            <CardContent className="space-y-4">
              <div className="flex items-center text-sm text-muted-foreground">
                <MapPin className="h-4 w-4 mr-1" />
                {getVenueText(event.venue)}
              </div>

              <div className="grid grid-cols-3 gap-4 text-center">
                <div>
                  <div className="text-2xl font-bold text-blue-600">{event.ticketsSold || 0}</div>
                  <div className="text-xs text-muted-foreground">Tickets Sold</div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-green-600">
                    {formatCurrency(event.revenue || 0)}
                  </div>
                  <div className="text-xs text-muted-foreground">Revenue</div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-purple-600">{event.capacity || 0}</div>
                  <div className="text-xs text-muted-foreground">/ capacity</div>
                </div>
              </div>

              <div className="flex gap-2">
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="flex-1"
                  onClick={() => navigate(`/admin/events/${event.id}`)}
                >
                  <Eye className="h-4 w-4 mr-1" />
                  View
                </Button>
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="flex-1"
                  onClick={() => navigate(`/admin/events/${event.id}`)}
                >
                  <Edit className="h-4 w-4 mr-1" />
                  Edit
                </Button>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {filteredEvents.length === 0 && !loading && (
        <div className="text-center py-12">
          <Calendar className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-medium mb-2">No events found</h3>
          <p className="text-muted-foreground mb-4">
            {searchTerm || statusFilter !== 'all' 
              ? 'Try adjusting your search or filters'
              : 'Get started by creating your first event'
            }
          </p>
          {(!searchTerm && statusFilter === 'all') && (
            <Button onClick={() => navigate('/admin/events/create')}>
              <Plus className="h-4 w-4 mr-2" />
              Create Event
            </Button>
          )}
        </div>
      )}
    </div>
  )
}

export default AdminEventsPage
