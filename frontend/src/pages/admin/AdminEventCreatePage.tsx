import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Trash2 } from 'lucide-react'

interface Category {
  id: string // UUID from backend
  name: string
  slug: string
  description: string
  icon?: string
  color?: string
  display_order?: number
  is_active?: boolean
}

const AdminEventCreatePage: React.FC = () => {
  const navigate = useNavigate()
  
  // Add categories state
  const [categories, setCategories] = useState<Category[]>([])
  const [loadingCategories, setLoadingCategories] = useState(true)
  
  // Get current date in proper format for datetime-local input
  const getCurrentDateTime = () => {
    // Set to October 18, 2025 8:00 PM as per UAT script
    return '2025-10-18T20:00'
  }
  
  // Get end date (4 hours after start)
  const getDefaultEndDateTime = () => {
    // Set to October 19, 2025 12:00 AM (4 hours after 8 PM)
    return '2025-10-19T00:00'
  }
  
  const [eventData, setEventData] = useState({
    title: 'Wizkid Live in Lagos',
    description: 'The Lagos grand finale of the Wizkid tour featuring special guest appearances and exclusive performances. Experience the biggest Afrobeats concert of the year at the iconic Eko Atlantic City.',
    category: 'Music',
    startDate: getCurrentDateTime(),
    endDate: getDefaultEndDateTime(),
    venueName: 'Eko Atlantic City',
    venueAddress: 'Eko Atlantic City, Victoria Island',
    venueCity: 'Lagos',
    venueState: 'lagos',
    venueCountry: 'Nigeria',
    venueCapacity: '50000',
    ticketTiers: [
      {
        name: 'VIP Diamond',
        description: 'Premium VIP experience with front row access, complimentary drinks, and meet & greet opportunity',
        price: '75000',
        quantity: '200',
        maxPerOrder: 2
      },
      {
        name: 'Golden Circle',
        description: 'Premium standing area with excellent stage view and exclusive bar access',
        price: '30000',
        quantity: '1500',
        maxPerOrder: 4
      }
    ],
    enableMomo: true,
    enablePaystack: true
  })

  // Fetch categories on component mount
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/categories')
        const data = await response.json()
        console.log('Categories API response:', data)
        if (data.success && data.data) {
          setCategories(data.data)
        }
      } catch (error) {
        console.error('Error fetching categories:', error)
      } finally {
        setLoadingCategories(false)
      }
    }

    fetchCategories()
  }, [])

  const handleInputChange = (field: string, value: string | number) => {
    setEventData(prev => {
      const updated = {
        ...prev,
        [field]: value
      }
      
      // Auto-update end date when start date changes (4 hours later)
      if (field === 'startDate' && typeof value === 'string' && value) {
        try {
          const startDate = new Date(value)
          if (!isNaN(startDate.getTime())) {
            const endDate = new Date(startDate)
            endDate.setHours(endDate.getHours() + 4) // 4 hours duration
            updated.endDate = endDate.toISOString().slice(0, 16)
          }
        } catch (error) {
          console.error('Error updating end date:', error)
        }
      }
      
      return updated
    })
  }

  // Multiple ticket tier management functions
  const addTicketTier = () => {
    setEventData(prev => ({
      ...prev,
      ticketTiers: [
        ...prev.ticketTiers,
        {
          name: '',
          description: '',
          price: '',
          quantity: '',
          maxPerOrder: 10
        }
      ]
    }))
  }

  const removeTicketTier = (index: number) => {
    if (eventData.ticketTiers.length > 1) {
      setEventData(prev => ({
        ...prev,
        ticketTiers: prev.ticketTiers.filter((_, i) => i !== index)
      }))
    }
  }

  const updateTicketTier = (index: number, field: string, value: string | number) => {
    setEventData(prev => ({
      ...prev,
      ticketTiers: prev.ticketTiers.map((tier, i) => 
        i === index ? { ...tier, [field]: value } : tier
      )
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      const adminToken = localStorage.getItem('adminToken')
      
      // Validate required fields
      if (!eventData.title.trim()) {
        alert('Please enter an event title')
        return
      }
      
      if (!eventData.category) {
        alert('Please select a category')
        return
      }
      
      if (!eventData.startDate) {
        alert('Please select a start date and time')
        return
      }
      
      if (!eventData.venueName.trim()) {
        alert('Please enter a venue name')
        return
      }
      
      if (!eventData.venueCity.trim()) {
        alert('Please enter a city')
        return
      }

      // Validate ticket tiers
      const validTiers = eventData.ticketTiers.filter(tier => 
        tier.name.trim() && tier.price && parseFloat(tier.price) > 0
      )
      
      if (validTiers.length === 0) {
        alert('Please add at least one valid ticket tier with name and price')
        return
      }

      // Create API payload matching backend expectations
      const createEventData = {
        organizer_id: '0ef31b06-6f14-4f7b-8b33-e62320f0b9bc', // Default organizer ID
        name: eventData.title,
        slug: eventData.title.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, ''),
        description: eventData.description,
        event_date: new Date(eventData.startDate).toISOString(),
        venue_name: eventData.venueName,
        venue_address: eventData.venueAddress || eventData.venueName,
        venue_city: eventData.venueCity,
        venue_state: eventData.venueState || null,
        venue_country: 'Nigeria',
        venue_capacity: parseInt(eventData.capacity) || null,
        ticket_tiers: validTiers.map(tier => ({
          name: tier.name.trim(),
          price: parseFloat(tier.price),
          quantity: parseInt(tier.quantity) || 100,
          max_per_order: tier.maxPerOrder || 10,
          description: tier.description || ''
        })),
        enable_momo: eventData.enableMomo,
        enable_paystack: eventData.enablePaystack
      }

      // Debug: Log the payload being sent
      console.log('=== CORRECTED API PAYLOAD ===')
      console.log(JSON.stringify(createEventData, null, 2))
      console.log('=============================')

      // Use the correct API endpoint
      const response = await fetch('http://localhost:8080/v1/admin/events', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(createEventData)
      })

      if (response.ok) {
        const result = await response.json()
        console.log('Event created successfully:', result)
        alert('Event created successfully!')
        navigate('/admin/dashboard')
      } else {
        const error = await response.text()
        console.error('Failed to create event:', error)
        alert('Failed to create event: ' + error)
      }
    } catch (error) {
      console.error('Error creating event:', error)
      alert('An error occurred while creating the event.')
    }
  }

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="mb-6">
        <button 
          onClick={() => navigate('/admin/dashboard')}
          className="mb-4 px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
        >
          ← Back to Dashboard
        </button>
        <h1 className="text-3xl font-bold">Create New Event</h1>
        <p className="text-gray-600">Set up a new event with ticket tiers and venue details</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Information */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-4">Basic Information</h2>
          
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium mb-2">Event Title *</label>
              <input
                type="text"
                required
                className="w-full p-2 border rounded-md"
                placeholder="Enter event title"
                value={eventData.title}
                onChange={(e) => handleInputChange('title', e.target.value)}
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium mb-2">Category *</label>
              <select 
                required
                className="w-full p-2 border rounded-md"
                value={eventData.category}
                onChange={(e) => handleInputChange('category', e.target.value)}
                disabled={loadingCategories}
              >
                <option value="">
                  {loadingCategories ? 'Loading categories...' : 'Select category'}
                </option>
                {categories.map(category => (
                  <option key={category.id} value={category.name}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="mt-4">
            <label className="block text-sm font-medium mb-2">Description</label>
            <textarea
              className="w-full p-2 border rounded-md"
              rows={4}
              placeholder="Describe your event..."
              value={eventData.description}
              onChange={(e) => handleInputChange('description', e.target.value)}
            />
          </div>

          <div className="grid gap-4 md:grid-cols-2 mt-4">
            <div>
              <label className="block text-sm font-medium mb-2">Start Date & Time *</label>
              <input
                type="datetime-local"
                required
                min={new Date().toISOString().slice(0, 16)}
                className="w-full p-2 border rounded-md"
                value={eventData.startDate}
                onChange={(e) => handleInputChange('startDate', e.target.value)}
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium mb-2">End Date & Time</label>
              <input
                type="datetime-local"
                min={eventData.startDate || new Date().toISOString().slice(0, 16)}
                className="w-full p-2 border rounded-md"
                value={eventData.endDate}
                onChange={(e) => handleInputChange('endDate', e.target.value)}
              />
            </div>
          </div>
        </div>

        {/* Venue Information */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-4">Venue Information</h2>
          
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium mb-2">Venue Name *</label>
              <input
                type="text"
                required
                className="w-full p-2 border rounded-md"
                placeholder="Enter venue name"
                value={eventData.venueName}
                onChange={(e) => handleInputChange('venueName', e.target.value)}
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium mb-2">Capacity</label>
              <input
                type="number"
                className="w-full p-2 border rounded-md"
                placeholder="Maximum attendees"
                value={eventData.venueCapacity}
                onChange={(e) => handleInputChange('venueCapacity', e.target.value)}
              />
            </div>
          </div>

          <div className="mt-4">
            <label className="block text-sm font-medium mb-2">Address</label>
            <input
              type="text"
              className="w-full p-2 border rounded-md"
              placeholder="Street address"
              value={eventData.venueAddress}
              onChange={(e) => handleInputChange('venueAddress', e.target.value)}
            />
          </div>

          <div className="grid gap-4 md:grid-cols-3 mt-4">
            <div>
              <label className="block text-sm font-medium mb-2">City *</label>
              <input
                type="text"
                required
                className="w-full p-2 border rounded-md"
                placeholder="City"
                value={eventData.venueCity}
                onChange={(e) => handleInputChange('venueCity', e.target.value)}
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium mb-2">State</label>
              <select 
                className="w-full p-2 border rounded-md"
                value={eventData.venueState}
                onChange={(e) => handleInputChange('venueState', e.target.value)}
              >
                <option value="">Select state</option>
                <option value="lagos">Lagos</option>
                <option value="abuja">FCT Abuja</option>
                <option value="rivers">Rivers</option>
                <option value="kano">Kano</option>
                <option value="oyo">Oyo</option>
                <option value="kaduna">Kaduna</option>
                <option value="other">Other</option>
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium mb-2">Country</label>
              <input
                type="text"
                className="w-full p-2 border rounded-md"
                value={eventData.venueCountry}
                onChange={(e) => handleInputChange('venueCountry', e.target.value)}
              />
            </div>
          </div>
        </div>

        {/* Multiple Ticket Tiers */}
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Ticket Tiers</h2>
            <button
              type="button"
              onClick={addTicketTier}
              className="flex items-center px-3 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              <Plus className="h-4 w-4 mr-1" />
              Add Tier
            </button>
          </div>
          
          {eventData.ticketTiers.map((tier, index) => (
            <div key={index} className="border rounded-lg p-4 mb-4 relative">
              {eventData.ticketTiers.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeTicketTier(index)}
                  className="absolute top-2 right-2 p-1 text-red-600 hover:bg-red-100 rounded"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              )}
              
              <h3 className="text-lg font-medium mb-3">Tier {index + 1}</h3>
              
              <div className="grid gap-4 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium mb-2">Tier Name *</label>
                  <input
                    type="text"
                    required
                    className="w-full p-2 border rounded-md"
                    placeholder="e.g., VIP, Regular, Student"
                    value={tier.name}
                    onChange={(e) => updateTicketTier(index, 'name', e.target.value)}
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-2">Price (NGN) *</label>
                  <input
                    type="number"
                    required
                    className="w-full p-2 border rounded-md"
                    placeholder="0"
                    value={tier.price}
                    onChange={(e) => updateTicketTier(index, 'price', e.target.value)}
                  />
                </div>
              </div>

              <div className="mt-4">
                <label className="block text-sm font-medium mb-2">Description</label>
                <textarea
                  className="w-full p-2 border rounded-md"
                  rows={3}
                  placeholder="Describe what's included in this tier..."
                  value={tier.description}
                  onChange={(e) => updateTicketTier(index, 'description', e.target.value)}
                />
              </div>

              <div className="grid gap-4 md:grid-cols-2 mt-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Quantity Available</label>
                  <input
                    type="number"
                    className="w-full p-2 border rounded-md"
                    placeholder="100"
                    value={tier.quantity}
                    onChange={(e) => updateTicketTier(index, 'quantity', e.target.value)}
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-2">Max Per Order</label>
                  <input
                    type="number"
                    className="w-full p-2 border rounded-md"
                    placeholder="10"
                    value={tier.maxPerOrder}
                    onChange={(e) => updateTicketTier(index, 'maxPerOrder', parseInt(e.target.value) || 10)}
                  />
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Payment Methods */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-4">Payment Methods</h2>
          <p className="text-sm text-gray-600 mb-4">Select which payment methods customers can use for this event</p>
          
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 border rounded-lg">
              <div className="flex-1">
                <h3 className="font-medium">MoMo PSB Payment</h3>
                <p className="text-sm text-gray-600">Enable mobile money payments via MoMo PSB</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  className="sr-only peer"
                  checked={eventData.enableMomo}
                  onChange={(e) => handleInputChange('enableMomo', e.target.checked)}
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            <div className="flex items-center justify-between p-4 border rounded-lg">
              <div className="flex-1">
                <h3 className="font-medium">Paystack Payment</h3>
                <p className="text-sm text-gray-600">Enable card and bank payments via Paystack</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  className="sr-only peer"
                  checked={eventData.enablePaystack}
                  onChange={(e) => handleInputChange('enablePaystack', e.target.checked)}
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>
          
          {!eventData.enableMomo && !eventData.enablePaystack && (
            <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
              <p className="text-sm text-yellow-800">⚠️ Warning: At least one payment method should be enabled for customers to purchase tickets.</p>
            </div>
          )}
        </div>

        {/* Submit Buttons */}
        <div className="flex justify-end space-x-4">
          <button
            type="button"
            onClick={() => navigate('/admin/dashboard')}
            className="px-6 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Create Event
          </button>
        </div>
      </form>
    </div>
  )
}

export default AdminEventCreatePage
