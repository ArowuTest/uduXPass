import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { eventsAPI } from '../services/api';
import { Event, PaginatedResponse, ApiResponse } from '../types/api';
import { formatCurrency, formatDate, debounce } from '../lib/utils';
import { 
  Calendar, 
  MapPin, 
  Star, 
  Search,
  Filter,
  ArrowRight,
  SlidersHorizontal,
  Ticket
} from 'lucide-react';
import { motion } from 'framer-motion';

interface FilterState {
  search: string;
  city: string;
  sort: string;
}

interface EventsPageState {
  events: Event[];
  isLoading: boolean;
  error: string | null;
  filters: FilterState;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

const EventsPage: React.FC = () => {
  const [state, setState] = useState<EventsPageState>({
    events: [],
    isLoading: true,
    error: null,
    filters: {
      search: '',
      city: '',
      sort: 'date'
    },
    pagination: {
      page: 1,
      limit: 12,
      total: 0,
      totalPages: 0
    }
  });

  // Load events from API
  const loadEvents = useCallback(async (filters?: FilterState, page?: number, limit?: number) => {
    console.log('[EventsPage] Loading events...');
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      
      const currentFilters = filters || state.filters;
      const currentPage = page || state.pagination.page;
      const currentLimit = limit || state.pagination.limit;
      
      const response: ApiResponse<PaginatedResponse<Event>> = await eventsAPI.getEvents({
        page: currentPage,
        limit: currentLimit,
        search: currentFilters.search || undefined,
        city: currentFilters.city || undefined,
      });
      
      if (response.success && response.data) {
        // Handle multiple response formats
        let eventsArray = [];
        if (Array.isArray(response.data.data)) {
          eventsArray = response.data.data;
        } else if (response.data.data?.events) {
          eventsArray = response.data.data.events;
        } else if (response.data.events) {
          eventsArray = response.data.events;
        }
        
        setState(prev => ({
          ...prev,
          events: eventsArray,
          pagination: {
            ...prev.pagination,
            total: response.data.pagination?.total || response.data.data?.pagination?.total || 0,
            totalPages: response.data.pagination?.total_pages || response.data.data?.pagination?.total_pages || 0,
          },
          isLoading: false
        }));
      } else {
        setState(prev => ({
          ...prev,
          events: [],
          isLoading: false,
          error: response.error || 'Failed to load events'
        }));
      }
    } catch (error) {
      console.error('Failed to load events:', error);
      setState(prev => ({
        ...prev,
        events: [],
        isLoading: false,
        error: error instanceof Error ? error.message : 'An unexpected error occurred'
      }));
    }
  }, [state.filters, state.pagination.page, state.pagination.limit]);

  // Trigger loads when filters change
  useEffect(() => {
    loadEvents(state.filters, state.pagination.page, state.pagination.limit);
  }, [state.pagination.page, state.pagination.limit, state.filters.search, state.filters.city]);

  // Initial load on mount
  useEffect(() => {
    console.log('[EventsPage] Component mounted, loading events...');
    loadEvents();
  }, []);

  // Debounced search handler
  const handleSearchChange = useCallback(
    debounce((value: string) => {
      setState(prev => ({
        ...prev,
        filters: { ...prev.filters, search: value },
        pagination: { ...prev.pagination, page: 1 } // Reset to page 1 on search
      }));
    }, 500),
    []
  );

  const handleCityChange = (city: string) => {
    setState(prev => ({
      ...prev,
      filters: { ...prev.filters, city },
      pagination: { ...prev.pagination, page: 1 }
    }));
  };

  const handlePageChange = (newPage: number) => {
    setState(prev => ({
      ...prev,
      pagination: { ...prev.pagination, page: newPage }
    }));
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const getLowestPrice = (event: Event): number => {
    if (!event.ticket_tiers || event.ticket_tiers.length === 0) return 0;
    return Math.min(...event.ticket_tiers.map(tier => tier.price));
  };

  const getAvailableTickets = (event: Event): number => {
    if (!event.ticket_tiers || event.ticket_tiers.length === 0) return 0;
    return event.ticket_tiers.reduce((sum, tier) => sum + (tier.quantity - tier.quantity_sold - tier.quantity_reserved), 0);
  };

  // Extract unique cities from events
  const cities = Array.from(new Set(state.events.map(e => e.venue?.city).filter(Boolean)));

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-secondary/10">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-purple-600 via-blue-600 to-indigo-600 text-white py-16">
        <div className="container mx-auto px-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="text-center"
          >
            <h1 className="text-5xl font-bold mb-4">Discover Amazing Events</h1>
            <p className="text-xl text-purple-100 mb-8">
              Find and book tickets to the hottest concerts, festivals, and shows
            </p>
            
            {/* Search Bar */}
            <div className="max-w-2xl mx-auto">
              <div className="relative">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <Input
                  type="text"
                  placeholder="Search events, artists, or venues..."
                  className="pl-12 pr-4 py-6 text-lg bg-white text-gray-900 rounded-full shadow-lg"
                  onChange={(e) => handleSearchChange(e.target.value)}
                />
              </div>
            </div>
          </motion.div>
        </div>
      </div>

      {/* Filters and Events */}
      <div className="container mx-auto px-4 py-12">
        {/* Filters */}
        <div className="flex flex-wrap gap-4 mb-8">
          <div className="flex items-center gap-2">
            <SlidersHorizontal className="h-5 w-5 text-gray-600" />
            <span className="font-semibold text-gray-700">Filters:</span>
          </div>
          
          {/* City Filter */}
          <select
            value={state.filters.city}
            onChange={(e) => handleCityChange(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
          >
            <option value="">All Cities</option>
            {cities.map(city => (
              <option key={city} value={city}>{city}</option>
            ))}
          </select>

          {/* Active Filters Display */}
          {(state.filters.search || state.filters.city) && (
            <div className="flex items-center gap-2">
              {state.filters.search && (
                <Badge variant="secondary" className="px-3 py-1">
                  Search: {state.filters.search}
                </Badge>
              )}
              {state.filters.city && (
                <Badge variant="secondary" className="px-3 py-1">
                  City: {state.filters.city}
                </Badge>
              )}
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setState(prev => ({
                  ...prev,
                  filters: { search: '', city: '', sort: 'date' },
                  pagination: { ...prev.pagination, page: 1 }
                }))}
              >
                Clear All
              </Button>
            </div>
          )}
        </div>

        {/* Results Count */}
        <div className="mb-6">
          <p className="text-gray-600">
            {state.isLoading ? 'Loading...' : `${state.pagination.total} events found`}
          </p>
        </div>

        {/* Loading State */}
        {state.isLoading && (
          <div className="flex justify-center items-center py-20">
            <LoadingSpinner size="lg" />
          </div>
        )}

        {/* Error State */}
        {state.error && !state.isLoading && (
          <div className="text-center py-20">
            <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md mx-auto">
              <p className="text-red-600 font-semibold mb-2">Error Loading Events</p>
              <p className="text-red-500 text-sm mb-4">{state.error}</p>
              <Button onClick={loadEvents} variant="outline">
                Try Again
              </Button>
            </div>
          </div>
        )}

        {/* Empty State */}
        {!state.isLoading && !state.error && state.events.length === 0 && (
          <div className="text-center py-20">
            <Ticket className="h-16 w-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-700 mb-2">No Events Found</h3>
            <p className="text-gray-500 mb-4">
              {state.filters.search || state.filters.city
                ? 'Try adjusting your filters to see more events'
                : 'Check back soon for upcoming events'}
            </p>
            {(state.filters.search || state.filters.city) && (
              <Button
                onClick={() => setState(prev => ({
                  ...prev,
                  filters: { search: '', city: '', sort: 'date' },
                  pagination: { ...prev.pagination, page: 1 }
                }))}
              >
                Clear Filters
              </Button>
            )}
          </div>
        )}

        {/* Events Grid */}
        {!state.isLoading && !state.error && state.events.length > 0 && (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
              {state.events.map((event, index) => (
                <motion.div
                  key={event.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.4, delay: index * 0.1 }}
                >
                  <Card className="h-full hover:shadow-xl transition-shadow duration-300 overflow-hidden group">
                    {/* Event Image */}
                    <div className="relative h-48 bg-gradient-to-br from-purple-500 to-blue-500 overflow-hidden">
                      {event.banner_image_url ? (
                        <img 
                          src={event.banner_image_url} 
                          alt={event.name}
                          className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                        />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center">
                          <Ticket className="h-16 w-16 text-white opacity-50" />
                        </div>
                      )}
                      <div className="absolute top-4 right-4">
                        <Badge className="bg-white text-purple-600 font-semibold">
                          {event.status === 'published' ? 'On Sale' : event.status}
                        </Badge>
                      </div>
                    </div>

                    <CardHeader>
                      <CardTitle className="text-xl font-bold line-clamp-2 group-hover:text-purple-600 transition-colors">
                        {event.name}
                      </CardTitle>
                    </CardHeader>

                    <CardContent className="space-y-3">
                      <div className="flex items-center text-gray-600">
                        <Calendar className="h-4 w-4 mr-2 flex-shrink-0" />
                        <span className="text-sm">{formatDate(event.event_date)}</span>
                      </div>
                      
                      <div className="flex items-start text-gray-600">
                        <MapPin className="h-4 w-4 mr-2 flex-shrink-0 mt-0.5" />
                        <div className="text-sm">
                          <p className="font-medium">{event.venue?.name}</p>
                          <p className="text-gray-500">{event.venue?.city}, {event.venue?.state}</p>
                        </div>
                      </div>

                      <div className="flex items-center justify-between pt-2 border-t">
                        <div>
                          <p className="text-xs text-gray-500">From</p>
                          <p className="text-lg font-bold text-purple-600">
                            {formatCurrency(getLowestPrice(event))}
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="text-xs text-gray-500">Available</p>
                          <p className="text-sm font-semibold text-gray-700">
                            {getAvailableTickets(event)} tickets
                          </p>
                        </div>
                      </div>
                    </CardContent>

                    <CardFooter>
                      <Link to={`/events/${event.id}`} className="w-full">
                        <Button className="w-full group-hover:bg-purple-600 transition-colors">
                          View Details
                          <ArrowRight className="ml-2 h-4 w-4" />
                        </Button>
                      </Link>
                    </CardFooter>
                  </Card>
                </motion.div>
              ))}
            </div>

            {/* Pagination */}
            {state.pagination.totalPages > 1 && (
              <div className="flex justify-center items-center gap-2">
                <Button
                  variant="outline"
                  onClick={() => handlePageChange(state.pagination.page - 1)}
                  disabled={state.pagination.page === 1}
                >
                  Previous
                </Button>
                
                <div className="flex gap-2">
                  {Array.from({ length: state.pagination.totalPages }, (_, i) => i + 1).map(page => (
                    <Button
                      key={page}
                      variant={page === state.pagination.page ? 'default' : 'outline'}
                      onClick={() => handlePageChange(page)}
                      className="min-w-[40px]"
                    >
                      {page}
                    </Button>
                  ))}
                </div>

                <Button
                  variant="outline"
                  onClick={() => handlePageChange(state.pagination.page + 1)}
                  disabled={state.pagination.page === state.pagination.totalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default EventsPage;
