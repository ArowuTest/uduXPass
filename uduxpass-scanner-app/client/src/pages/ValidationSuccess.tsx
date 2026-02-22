/**
 * Design Philosophy: Professional Event Tech - Full-screen green success state
 * Layout: Gradient background, checkmark animation, ticket details card
 * Colors: Vibrant green (#10B981) for positive reinforcement
 */

import { useEffect } from 'react';
import { useLocation } from 'wouter';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { CheckCircle, User, Ticket, Calendar, MapPin } from 'lucide-react';

export default function ValidationSuccess() {
  const [, setLocation] = useLocation();
  const [location] = useLocation();
  
  // Get ticket data from navigation state
  const state = (window.history.state as any)?.state;
  const ticket = state?.ticket;
  const message = state?.message || 'Valid Ticket';

  useEffect(() => {
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate([100, 50, 100]);
    }

    // Auto-redirect after 5 seconds
    const timer = setTimeout(() => {
      handleScanNext();
    }, 5000);

    return () => clearTimeout(timer);
  }, []);

  const handleScanNext = () => {
    setLocation('/scan');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-500 to-green-600 flex flex-col items-center justify-center p-6 text-white">
      {/* Success Icon with Animation */}
      <div className="mb-6 animate-in zoom-in duration-500">
        <div className="relative">
          <div className="absolute inset-0 bg-white/20 rounded-full animate-ping" />
          <CheckCircle className="h-24 w-24 relative z-10" strokeWidth={2} />
        </div>
      </div>

      {/* Heading */}
      <h1 className="text-3xl font-bold mb-2 animate-in fade-in slide-in-from-bottom-4 duration-500 delay-100">
        Valid Ticket
      </h1>
      <p className="text-green-100 mb-8 animate-in fade-in slide-in-from-bottom-4 duration-500 delay-200">
        {message}
      </p>

      {/* Ticket Details Card */}
      {ticket && (
        <Card className="w-full max-w-md bg-white text-foreground shadow-2xl animate-in fade-in slide-in-from-bottom-4 duration-500 delay-300">
          <CardContent className="p-6 space-y-4">
            {/* Attendee Name */}
            {ticket.user && (
              <div className="flex items-start gap-3">
                <User className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm text-muted-foreground">Attendee Name</p>
                  <p className="font-semibold text-lg">{ticket.user.name}</p>
                </div>
              </div>
            )}

            {/* Ticket Type */}
            {ticket.ticket_type && (
              <div className="flex items-start gap-3">
                <Ticket className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm text-muted-foreground">Ticket Type</p>
                  <p className="font-semibold">{ticket.ticket_type.name}</p>
                </div>
              </div>
            )}

            {/* Event */}
            {ticket.event && (
              <div className="flex items-start gap-3">
                <Calendar className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm text-muted-foreground">Event</p>
                  <p className="font-semibold">{ticket.event.name}</p>
                </div>
              </div>
            )}

            {/* Location */}
            {ticket.event?.location && (
              <div className="flex items-start gap-3">
                <MapPin className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm text-muted-foreground">Location</p>
                  <p className="font-medium">{ticket.event.location}</p>
                </div>
              </div>
            )}

            {/* Scan Timestamp */}
            <div className="pt-4 border-t border-border text-center">
              <p className="text-sm text-muted-foreground">
                Scanned at {new Date().toLocaleTimeString()}
              </p>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Action Button */}
      <Button
        onClick={handleScanNext}
        className="mt-8 w-full max-w-md h-14 text-base font-semibold bg-white text-green-600 hover:bg-green-50 animate-in fade-in slide-in-from-bottom-4 duration-500 delay-500"
      >
        Scan Next Ticket
      </Button>

      {/* Auto-redirect notice */}
      <p className="mt-4 text-sm text-green-100 animate-in fade-in duration-500 delay-700">
        Redirecting in 5 seconds...
      </p>
    </div>
  );
}
