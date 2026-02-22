import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import LoadingSpinner from '@/components/ui/LoadingSpinner';
import { eventsAPI } from '@/services/api';
import { Event, PaginatedResponse, ApiResponse } from '@/types/api';
import { formatCurrency, formatDate } from '@/lib/utils';
import { 
  Calendar, 
  MapPin, 
  Users, 
  Ticket, 
  Shield, 
  Smartphone,
  ArrowRight,
  Play,
  LucideIcon
} from 'lucide-react';

interface Feature {
  icon: LucideIcon;
  title: string;
  description: string;
}

interface HomePageState {
  featuredEvents: Event[];
  isLoading: boolean;
  error: string | null;
}

const HomePage: React.FC = () => {
  const [state, setState] = useState<HomePageState>({
    featuredEvents: [],
    isLoading: true,
    error: null
  });

  useEffect(() => {
    loadFeaturedEvents();
  }, []);

  const loadFeaturedEvents = async (): Promise<void> => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      
      const response: ApiResponse<PaginatedResponse<Event>> = await eventsAPI.getEvents({ 
        limit: 6,
        status: 'published'
      });
      
      if (response.success && response.data?.data) {
        setState(prev => ({
          ...prev,
          featuredEvents: response.data!.data,
          isLoading: false
        }));
      } else {
        setState(prev => ({
          ...prev,
          featuredEvents: [],
          isLoading: false,
          error: response.error || 'Failed to load events'
        }));
      }
    } catch (error) {
      console.error('Failed to load featured events:', error);
      setState(prev => ({
        ...prev,
        featuredEvents: [],
        isLoading: false,
        error: error instanceof Error ? error.message : 'An unexpected error occurred'
      }));
    }
  };

  const getLowestPrice = (event: Event): number => {
    if (!event.ticket_tiers || event.ticket_tiers.length === 0) return 0;
    return Math.min(...event.ticket_tiers.map(tier => tier.price));
  };

  const features: Feature[] = [
    {
      icon: Shield,
      title: 'Secure Payments',
      description: 'Pay safely with MoMo or Paystack integration'
    },
    {
      icon: Smartphone,
      title: 'Mobile First',
      description: 'Seamless experience on all your devices'
    },
    {
      icon: Ticket,
      title: 'Digital Tickets',
      description: 'QR code tickets delivered instantly'
    },
    {
      icon: Users,
      title: 'Trusted Platform',
      description: 'Join thousands of satisfied customers'
    }
  ];

  const { featuredEvents, isLoading, error } = state;

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900 text-white overflow-hidden">
        <div className="absolute inset-0 bg-black/20"></div>
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 lg:py-32">
          <div className="text-center max-w-4xl mx-auto">
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold mb-6 animate-fade-in">
              Experience
              <span className="block bg-gradient-to-r from-yellow-400 to-orange-500 bg-clip-text text-transparent">
                Unforgettable
              </span>
              Events
            </h1>
            <p className="text-xl md:text-2xl text-gray-300 mb-8 leading-relaxed animate-fade-in-delay">
              Discover amazing events, secure your tickets instantly, and create memories that last a lifetime
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center animate-fade-in-delay-2">
              <Button asChild size="lg" className="bg-gradient-to-r from-yellow-400 to-orange-500 hover:from-yellow-500 hover:to-orange-600 text-black font-semibold transform hover:scale-105 transition-all duration-200">
                <Link to="/events">
                  Browse Events <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="border-white text-white hover:bg-white hover:text-black transform hover:scale-105 transition-all duration-200">
                <Play className="mr-2 h-5 w-5" />
                Watch Demo
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Why Choose uduXPass?
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              The most trusted platform for event ticketing in Nigeria
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <div
                key={feature.title}
                className="transform hover:scale-105 transition-all duration-300"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <Card className="text-center h-full hover:shadow-lg transition-shadow">
                  <CardHeader>
                    <div className="mx-auto w-16 h-16 bg-gradient-to-br from-purple-500 to-blue-600 rounded-full flex items-center justify-center mb-4">
                      <feature.icon className="h-8 w-8 text-white" />
                    </div>
                    <CardTitle className="text-xl font-semibold">{feature.title}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-gray-600">{feature.description}</p>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Featured Events Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Featured Events
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Don't miss out on these amazing upcoming events
            </p>
          </div>

          {isLoading ? (
            <div className="flex justify-center">
              <LoadingSpinner size="lg" />
            </div>
          ) : error ? (
            <div className="text-center py-12">
              <div className="text-red-500 text-6xl mb-4">⚠️</div>
              <h3 className="text-xl font-semibold text-gray-600 mb-2">Failed to Load Events</h3>
              <p className="text-gray-500 mb-4">{error}</p>
              <Button onClick={loadFeaturedEvents} variant="outline">
                Try Again
              </Button>
            </div>
          ) : featuredEvents.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {featuredEvents.map((event, index) => (
                <div
                  key={event.id}
                  className="transform hover:scale-105 transition-all duration-300"
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  <Card className="overflow-hidden hover:shadow-xl transition-shadow group">
                    <div className="relative">
                      <img 
                        src={event.event_image_url || '/api/placeholder/400/250'} 
                        alt={event.name}
                        className="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-300"
                      />
                      <div className="absolute top-4 right-4">
                        <Badge variant="secondary" className="bg-white/90 text-gray-800">
                          {event.status}
                        </Badge>
                      </div>
                    </div>
                    <CardHeader>
                      <CardTitle className="line-clamp-2">{event.name}</CardTitle>
                      <div className="flex items-center text-sm text-gray-500 space-x-4">
                        <div className="flex items-center">
                          <Calendar className="h-4 w-4 mr-1" />
                          {formatDate(event.event_date)}
                        </div>
                        <div className="flex items-center">
                          <MapPin className="h-4 w-4 mr-1" />
                          {event.venue_name}
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <p className="text-gray-600 line-clamp-3 mb-4">{event.description}</p>
                      <div className="flex items-center justify-between">
                        <div className="text-2xl font-bold text-purple-600">
                          From {formatCurrency(getLowestPrice(event))}
                        </div>
                        <div className="flex items-center text-sm text-gray-500">
                          <Users className="h-4 w-4 mr-1" />
                          {event.ticket_tiers?.length || 0} tiers
                        </div>
                      </div>
                    </CardContent>
                    <CardFooter>
                      <Button asChild className="w-full">
                        <Link to={`/events/${event.id}`}>
                          View Details
                        </Link>
                      </Button>
                    </CardFooter>
                  </Card>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <Ticket className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <h3 className="text-xl font-semibold text-gray-600 mb-2">No Events Available</h3>
              <p className="text-gray-500">Check back soon for exciting events!</p>
            </div>
          )}

          <div className="text-center mt-12">
            <Button asChild size="lg" variant="outline">
              <Link to="/events">
                View All Events <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-purple-600 to-blue-600 text-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div>
            <h2 className="text-3xl md:text-4xl font-bold mb-6">
              Ready to Experience Amazing Events?
            </h2>
            <p className="text-xl mb-8 text-purple-100">
              Join thousands of event-goers who trust uduXPass for their ticketing needs
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button asChild size="lg" className="bg-white text-purple-600 hover:bg-gray-100 transform hover:scale-105 transition-all duration-200">
                <Link to="/events">
                  Browse Events
                </Link>
              </Button>
              <Button asChild variant="outline" size="lg" className="border-white text-white hover:bg-white hover:text-purple-600 transform hover:scale-105 transition-all duration-200">
                <Link to="/register">
                  Create Account
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default HomePage;

