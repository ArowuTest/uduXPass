import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useToast } from '@/components/ui/use-toast';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { useCart } from '../contexts/CartContext';
import { useAuth } from '../contexts/AuthContext';
import { eventsAPI } from '../services/api';
import { Event } from '../types/api';
import { formatCurrency, formatDate } from '../lib/utils';
import { 
  Calendar, 
  MapPin, 
  ArrowLeft,
  Plus,
  Minus,
  ShoppingCart,
  Clock,
  Users,
  Share2,
  Heart,
  Ticket,
  Info
} from 'lucide-react';
import { motion } from 'framer-motion';

interface TicketSelection {
  tierId: string;
  tierName: string;
  quantity: number;
  price: number;
}

const EventDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { addItem } = useCart();
  const { isAuthenticated } = useAuth();
  
  const [event, setEvent] = useState<Event | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [ticketSelections, setTicketSelections] = useState<{ [tierId: string]: number }>({});
  const [isFavorite, setIsFavorite] = useState<boolean>(false);

  useEffect(() => {
    if (id) {
      loadEvent();
    }
  }, [id]);

  const loadEvent = async (): Promise<void> => {
    if (!id) return;
    
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await eventsAPI.getById(id);
      
      if (response.success && response.data) {
        setEvent(response.data);
      } else {
        setError(response.error || 'Failed to load event');
      }
    } catch (error) {
      console.error('Error loading event:', error);
      setError('Failed to load event details');
    } finally {
      setIsLoading(false);
    }
  };

  const formatTime = (dateString: string): string => {
    return new Date(dateString).toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    });
  };

  const handleTicketQuantityChange = (tierId: string, quantity: number, maxQuantity: number) => {
    const newQuantity = Math.max(0, Math.min(quantity, Math.min(maxQuantity, 10)));
    setTicketSelections(prev => ({
      ...prev,
      [tierId]: newQuantity
    }));
  };

  const getTotalTickets = (): number => {
    return Object.values(ticketSelections).reduce((sum, qty) => sum + qty, 0);
  };

  const getTotalPrice = (): number => {
    if (!event?.ticket_tiers) return 0;
    
    return event.ticket_tiers.reduce((total, tier) => {
      const quantity = ticketSelections[tier.id] || 0;
      return total + (tier.price * quantity);
    }, 0);
  };

  const getAvailableQuantity = (tier: any): number => {
    // Backend returns: quota (total), sold (number sold)
    // Available = quota - sold
    const quota = tier.quota || 0;
    const sold = tier.sold || 0;
    return Math.max(0, quota - sold);
  };

  const handleAddToCart = () => {
    if (!event) return;
    
    if (!isAuthenticated) {
      toast({
        title: "Authentication Required",
        description: "Please log in to purchase tickets",
        variant: "destructive",
      });
      navigate('/login', { state: { from: `/events/${id}` } });
      return;
    }
    
    const selections: TicketSelection[] = event.ticket_tiers
      .filter(tier => ticketSelections[tier.id] > 0)
      .map(tier => ({
        tierId: tier.id,
        tierName: tier.name,
        quantity: ticketSelections[tier.id],
        price: tier.price
      }));
    
    if (selections.length === 0) {
      toast({
        title: "No Tickets Selected",
        description: "Please select at least one ticket",
        variant: "destructive",
      });
      return;
    }

    // Add to cart
    selections.forEach(selection => {
      const tier = event.ticket_tiers.find(t => t.id === selection.tierId);
      if (tier) {
        addItem(event.id, tier, selection.quantity);
      }
    });

    toast({
      title: "Added to Cart",
      description: `${getTotalTickets()} ticket(s) added to your cart`,
    });

    // Navigate to checkout
    navigate('/checkout');
  };

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: event?.name,
          text: event?.description,
          url: window.location.href,
        });
      } catch (error) {
        console.error('Error sharing:', error);
      }
    } else {
      // Fallback: Copy to clipboard
      navigator.clipboard.writeText(window.location.href);
      toast({
        title: "Link Copied",
        description: "Event link copied to clipboard",
      });
    }
  };

  const handleToggleFavorite = () => {
    setIsFavorite(!isFavorite);
    toast({
      title: isFavorite ? "Removed from Favorites" : "Added to Favorites",
      description: isFavorite ? "Event removed from your favorites" : "Event added to your favorites",
    });
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="flex flex-col items-center justify-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-gray-600">Loading event details...</p>
        </div>
      </div>
    );
  }

  if (error || !event) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-md mx-auto text-center">
          <div className="bg-red-50 border border-red-200 rounded-lg p-8">
            <Ticket className="h-16 w-16 text-red-300 mx-auto mb-4" />
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Event Not Found</h1>
            <p className="text-gray-600 mb-6">{error || 'The event you\'re looking for doesn\'t exist.'}</p>
            <Button onClick={() => navigate('/events')} className="w-full">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Events
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const totalAvailableTickets = event.ticket_tiers?.reduce((sum, tier) => sum + getAvailableQuantity(tier), 0) || 0;

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-secondary/10">
      {/* Hero Section with Event Banner */}
      <div className="relative h-96 bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-600 overflow-hidden">
        {event.banner_image_url && (
          <img 
            src={event.banner_image_url} 
            alt={event.name}
            className="absolute inset-0 w-full h-full object-cover opacity-40"
          />
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
        
        <div className="relative container mx-auto px-4 h-full flex flex-col justify-end pb-8">
          <Button 
            variant="ghost" 
            onClick={() => navigate('/events')}
            className="absolute top-4 left-4 text-white hover:bg-white/20"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>

          <div className="flex items-center gap-2 mb-4">
            <Badge className="bg-green-500 text-white">
              {event.status === 'published' ? 'On Sale' : event.status}
            </Badge>
            <Badge variant="outline" className="bg-white/20 text-white border-white/40">
              {totalAvailableTickets} tickets available
            </Badge>
          </div>

          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">{event.name}</h1>
          
          <div className="flex flex-wrap gap-6 text-white">
            <div className="flex items-center">
              <Calendar className="h-5 w-5 mr-2" />
              <div>
                <p className="font-semibold">{formatDate(event.event_date)}</p>
                <p className="text-sm text-white/80">{formatTime(event.event_date)}</p>
              </div>
            </div>
            
            <div className="flex items-center">
              <MapPin className="h-5 w-5 mr-2" />
              <div>
                <p className="font-semibold">{event.venue?.name}</p>
                <p className="text-sm text-white/80">{event.venue?.city}, {event.venue?.state}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-12">
        <div className="grid gap-8 lg:grid-cols-3">
          {/* Event Details */}
          <div className="lg:col-span-2 space-y-6">
            {/* About Event */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4 }}
            >
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Info className="h-5 w-5 mr-2" />
                    About This Event
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-gray-700 leading-relaxed whitespace-pre-line">
                    {event.description}
                  </p>
                </CardContent>
              </Card>
            </motion.div>

            {/* Venue Details */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: 0.1 }}
            >
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <MapPin className="h-5 w-5 mr-2" />
                    Venue Information
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div>
                    <p className="font-semibold text-lg">{event.venue?.name}</p>
                    <p className="text-gray-600">{event.venue?.address}</p>
                    <p className="text-gray-600">{event.venue?.city}, {event.venue?.state}, {event.venue?.country}</p>
                  </div>
                  {event.venue?.capacity && (
                    <div className="flex items-center text-gray-600">
                      <Users className="h-4 w-4 mr-2" />
                      <span>Capacity: {event.venue.capacity.toLocaleString()}</span>
                    </div>
                  )}
                </CardContent>
              </Card>
            </motion.div>

            {/* Event Stats */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: 0.2 }}
            >
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <Card>
                  <CardContent className="pt-6 text-center">
                    <Ticket className="h-8 w-8 text-purple-600 mx-auto mb-2" />
                    <p className="text-2xl font-bold">{totalAvailableTickets}</p>
                    <p className="text-sm text-gray-600">Available</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6 text-center">
                    <Users className="h-8 w-8 text-blue-600 mx-auto mb-2" />
                    <p className="text-2xl font-bold">{event.tickets_sold || 0}</p>
                    <p className="text-sm text-gray-600">Sold</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6 text-center">
                    <Calendar className="h-8 w-8 text-green-600 mx-auto mb-2" />
                    <p className="text-2xl font-bold">{event.ticket_tiers?.length || 0}</p>
                    <p className="text-sm text-gray-600">Ticket Tiers</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6 text-center">
                    <Clock className="h-8 w-8 text-orange-600 mx-auto mb-2" />
                    <p className="text-2xl font-bold">10</p>
                    <p className="text-sm text-gray-600">Min Hold</p>
                  </CardContent>
                </Card>
              </div>
            </motion.div>
          </div>

          {/* Ticket Selection Sidebar */}
          <div className="lg:col-span-1">
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.4 }}
              className="sticky top-8"
            >
              <Card className="shadow-xl">
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span>Select Tickets</span>
                    <div className="flex gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleToggleFavorite}
                        className={isFavorite ? 'text-red-500' : ''}
                      >
                        <Heart className={`h-5 w-5 ${isFavorite ? 'fill-current' : ''}`} />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleShare}
                      >
                        <Share2 className="h-5 w-5" />
                      </Button>
                    </div>
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {event.ticket_tiers && event.ticket_tiers.length > 0 ? (
                    <>
                      {event.ticket_tiers.map((tier) => {
                        const available = getAvailableQuantity(tier);
                        const isAvailable = available > 0;
                        
                        return (
                          <div 
                            key={tier.id} 
                            className={`border rounded-lg p-4 transition-all ${
                              isAvailable ? 'hover:border-purple-500 hover:shadow-md' : 'opacity-50'
                            }`}
                          >
                            <div className="flex justify-between items-start mb-3">
                              <div className="flex-1">
                                <h4 className="font-semibold text-lg">{tier.name}</h4>
                                {tier.description && (
                                  <p className="text-sm text-gray-600 mt-1">{tier.description}</p>
                                )}
                                <p className="text-sm text-gray-500 mt-2">
                                  {available} of {tier.quantity} available
                                </p>
                              </div>
                              <div className="text-right">
                                <p className="font-bold text-xl text-purple-600">
                                  {formatCurrency(tier.price)}
                                </p>
                              </div>
                            </div>
                            
                            {isAvailable ? (
                              <div className="flex items-center justify-between mt-3 pt-3 border-t">
                                <span className="text-sm text-gray-600">Quantity:</span>
                                <div className="flex items-center gap-2">
                                  <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => handleTicketQuantityChange(
                                      tier.id, 
                                      (ticketSelections[tier.id] || 0) - 1,
                                      available
                                    )}
                                    disabled={!ticketSelections[tier.id] || ticketSelections[tier.id] <= 0}
                                  >
                                    <Minus className="h-4 w-4" />
                                  </Button>
                                  
                                  <span className="w-12 text-center font-semibold text-lg">
                                    {ticketSelections[tier.id] || 0}
                                  </span>
                                  
                                  <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => handleTicketQuantityChange(
                                      tier.id, 
                                      (ticketSelections[tier.id] || 0) + 1,
                                      available
                                    )}
                                    disabled={(ticketSelections[tier.id] || 0) >= Math.min(available, 10)}
                                  >
                                    <Plus className="h-4 w-4" />
                                  </Button>
                                </div>
                              </div>
                            ) : (
                              <div className="mt-3 pt-3 border-t">
                                <Badge variant="destructive" className="w-full justify-center">
                                  Sold Out
                                </Badge>
                              </div>
                            )}
                          </div>
                        );
                      })}

                      {/* Cart Summary */}
                      {getTotalTickets() > 0 && (
                        <div className="border-t pt-4 space-y-4">
                          <div className="space-y-2">
                            <div className="flex justify-between text-sm">
                              <span className="text-gray-600">Subtotal ({getTotalTickets()} tickets)</span>
                              <span className="font-semibold">{formatCurrency(getTotalPrice())}</span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-gray-600">Service Fee</span>
                              <span className="font-semibold">{formatCurrency(getTotalPrice() * 0.05)}</span>
                            </div>
                            <div className="flex justify-between text-lg font-bold pt-2 border-t">
                              <span>Total</span>
                              <span className="text-purple-600">{formatCurrency(getTotalPrice() * 1.05)}</span>
                            </div>
                          </div>
                          
                          <Button 
                            className="w-full bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white font-semibold py-6 text-lg" 
                            onClick={handleAddToCart}
                          >
                            <ShoppingCart className="h-5 w-5 mr-2" />
                            Proceed to Checkout
                          </Button>

                          <p className="text-xs text-center text-gray-500">
                            Tickets will be held for 10 minutes during checkout
                          </p>
                        </div>
                      )}
                    </>
                  ) : (
                    <div className="text-center py-8">
                      <Ticket className="h-12 w-12 text-gray-300 mx-auto mb-3" />
                      <p className="text-gray-600">No tickets available at this time</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </motion.div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EventDetailsPage;
