/**
 * Design Philosophy: Professional Event Tech - Full-screen red error state
 * Layout: Gradient background, warning icon, error details card
 * Colors: Bold red (#EF4444) for clear alert signal
 */

import { useEffect } from 'react';
import { useLocation } from 'wouter';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { AlertTriangle, AlertCircle } from 'lucide-react';

export default function ValidationError() {
  const [, setLocation] = useLocation();
  
  // Get error data from navigation state
  const state = (window.history.state as any)?.state;
  const message = state?.message || 'Invalid Ticket';
  const ticket = state?.ticket;
  const errorType = state?.errorType || 'ALREADY_USED'; // ALREADY_USED (red), INVALID (yellow), NOT_CACHED (yellow)
  
  // Determine colors based on error type
  const isInvalidTicket = errorType === 'INVALID' || errorType === 'NOT_CACHED';
  const bgGradient = isInvalidTicket 
    ? 'from-yellow-500 to-yellow-600' 
    : 'from-red-500 to-red-600';
  const textColor = isInvalidTicket ? 'text-yellow-600' : 'text-red-600';
  const buttonTextColor = isInvalidTicket ? 'text-yellow-600' : 'text-red-600';
  const buttonHoverBg = isInvalidTicket ? 'hover:bg-yellow-50' : 'hover:bg-red-50';
  const heading = isInvalidTicket ? 'Invalid Ticket' : 'Ticket Already Used';

  useEffect(() => {
    // Error haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate([200, 100, 200]);
    }
  }, []);

  const handleScanNext = () => {
    setLocation('/scan');
  };

  const handleOverride = () => {
    // TODO: Implement override functionality (admin only)
    // For now, just show a toast
    alert('Override functionality requires admin permissions');
  };

  return (
    <div className={`min-h-screen bg-gradient-to-br ${bgGradient} flex flex-col items-center justify-center p-6 text-white`}>
      {/* Warning Icon with Animation */}
      <div className="mb-6 animate-in zoom-in duration-500">
        <div className="relative">
          <div className="absolute inset-0 bg-white/20 rounded-full animate-ping" />
          <AlertTriangle className="h-24 w-24 relative z-10" strokeWidth={2} />
        </div>
      </div>

      {/* Heading */}
      <h1 className="text-3xl font-bold mb-8 animate-in fade-in slide-in-from-bottom-4 duration-500 delay-100">
        {heading}
      </h1>

      {/* Error Details Card */}
      <Card className="w-full max-w-md bg-white text-foreground shadow-2xl animate-in fade-in slide-in-from-bottom-4 duration-500 delay-200">
        <CardContent className="p-6 space-y-4">
          {/* Error Reason */}
          <div className="flex items-start gap-3">
            <AlertCircle className={`h-5 w-5 ${textColor} mt-0.5 flex-shrink-0`} />
            <div>
              <p className="text-sm text-muted-foreground mb-1">Reason:</p>
              <p className={`font-semibold text-lg ${textColor}`}>{message}</p>
            </div>
          </div>

          {/* Additional Ticket Info */}
          {ticket && (
            <>
              {ticket.scanned_at && ticket.scanned_by && (
                <div className="pt-4 border-t border-border">
                  <p className="text-sm text-muted-foreground mb-1">Previously scanned:</p>
                  <p className="font-medium">
                    {new Date(ticket.scanned_at).toLocaleString()}
                  </p>
                  <p className="text-sm text-muted-foreground mt-1">
                    by Scanner #{ticket.scanned_by}
                  </p>
                </div>
              )}

              {ticket.ticket_code && (
                <div className="pt-4 border-t border-border">
                  <p className="text-sm text-muted-foreground mb-1">Ticket ID:</p>
                  <p className="font-mono text-sm">{ticket.ticket_code}</p>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      {/* Action Buttons */}
      <div className="w-full max-w-md space-y-3 mt-8">
        {/* Override Button (Admin Only) */}
        <Button
          onClick={handleOverride}
          variant="outline"
          className="w-full h-14 text-base font-semibold bg-transparent border-2 border-white text-white hover:bg-white/10 animate-in fade-in slide-in-from-bottom-4 duration-500 delay-300"
        >
          Override & Allow Entry
        </Button>

        {/* Scan Next Button */}
        <Button
          onClick={handleScanNext}
          className={`w-full h-14 text-base font-semibold bg-white ${buttonTextColor} ${buttonHoverBg} animate-in fade-in slide-in-from-bottom-4 duration-500 delay-400`}
        >
          Scan Next Ticket
        </Button>
      </div>
    </div>
  );
}
