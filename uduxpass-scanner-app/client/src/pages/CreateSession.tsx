/**
 * Design Philosophy: Professional Event Tech - Clean form, clear inputs
 * Layout: Simple form with event selector, location input, notes
 */

import { useState, useEffect } from 'react';
import { useLocation } from 'wouter';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { scannerApi, Event } from '@/lib/api';
import { ArrowLeft, Loader2, MapPin, Download } from 'lucide-react';
import { toast } from 'sonner';
import { offlineDB } from '@/lib/offlineDB';

export default function CreateSession() {
  const [, setLocation] = useLocation();
  const [events, setEvents] = useState<Event[]>([]);
  const [selectedEventId, setSelectedEventId] = useState('');
  const [location, setLocationValue] = useState('');
  const [notes, setNotes] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingEvents, setIsLoadingEvents] = useState(true);
  const [isCaching, setIsCaching] = useState(false);

  useEffect(() => {
    loadEvents();
  }, []);

  const loadEvents = async () => {
    try {
      const eventsData = await scannerApi.getEvents();
      setEvents(eventsData);
    } catch (error) {
      console.error('Failed to load events:', error);
      toast.error('Failed to load events');
    } finally {
      setIsLoadingEvents(false);
    }
  };

  const cacheTickets = async () => {
    if (!selectedEventId) {
      toast.error('Please select an event first');
      return;
    }

    setIsCaching(true);
    try {
      // Fetch all tickets for the event from API
      const response = await fetch(`/api/scanner/events/${selectedEventId}/tickets`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error('Failed to fetch tickets');
      }
      
      const data = await response.json();
      const tickets = data.data || [];
      
      // Cache tickets in IndexedDB
      await offlineDB.cacheTickets(selectedEventId, tickets);
      
      const stats = await offlineDB.getStats();
      toast.success(`Cached ${tickets.length} tickets for offline use`);
      console.log('Offline DB stats:', stats);
    } catch (error) {
      console.error('Failed to cache tickets:', error);
      toast.error('Failed to cache tickets for offline use');
    } finally {
      setIsCaching(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!selectedEventId) {
      toast.error('Please select an event');
      return;
    }

    setIsLoading(true);
    try {
      await scannerApi.createSession({
        event_id: selectedEventId,
      });

      toast.success('Scanning session created successfully!');
      setLocation('/dashboard');
    } catch (error: any) {
      console.error('Failed to create session:', error);
      toast.error(error.response?.data?.message || 'Failed to create session');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Top Bar */}
      <header className="bg-card border-b border-border px-4 py-3 flex items-center gap-3">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setLocation('/dashboard')}
          className="text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <h1 className="font-semibold text-lg">Create Session</h1>
      </header>

      {/* Form */}
      <main className="p-4">
        <form onSubmit={handleSubmit} className="space-y-6 max-w-md mx-auto">
          {/* Event Selector */}
          <div className="space-y-2">
            <Label htmlFor="event">Event *</Label>
            {isLoadingEvents ? (
              <div className="h-12 bg-muted rounded-md flex items-center justify-center">
                <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
              </div>
            ) : (
              <Select value={selectedEventId} onValueChange={setSelectedEventId}>
                <SelectTrigger id="event" className="h-12">
                  <SelectValue placeholder="Select Event" />
                </SelectTrigger>
                <SelectContent>
                  {events.length === 0 ? (
                    <div className="p-4 text-center text-sm text-muted-foreground">
                      No events available
                    </div>
                  ) : (
                    events.map((event) => (
                      <SelectItem key={event.id} value={event.id}>
                        {event.name}
                      </SelectItem>
                    ))
                  )}
                </SelectContent>
              </Select>
            )}
            
            {/* Cache Tickets Button */}
            {selectedEventId && (
              <Button
                type="button"
                variant="outline"
                onClick={cacheTickets}
                disabled={isCaching}
                className="w-full mt-2"
              >
                {isCaching ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Caching tickets...
                  </>
                ) : (
                  <>
                    <Download className="mr-2 h-4 w-4" />
                    Cache Tickets for Offline Use
                  </>
                )}
              </Button>
            )}
          </div>

          {/* Location Input */}
          <div className="space-y-2">
            <Label htmlFor="location">Location / Entrance *</Label>
            <div className="relative">
              <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
              <Input
                id="location"
                type="text"
                placeholder="Main Entrance"
                value={location}
                onChange={(e) => setLocationValue(e.target.value)}
                className="pl-10 h-12"
                disabled={isLoading}
              />
            </div>
          </div>

          {/* Notes Textarea */}
          <div className="space-y-2">
            <Label htmlFor="notes">Notes (Optional)</Label>
            <Textarea
              id="notes"
              placeholder="Add any notes about this scanning session..."
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="min-h-24 resize-none"
              disabled={isLoading}
            />
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            className="w-full h-14 text-base font-semibold"
            disabled={isLoading || isLoadingEvents || events.length === 0}
          >
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                Creating Session...
              </>
            ) : (
              'Start Scanning Session'
            )}
          </Button>
        </form>
      </main>
    </div>
  );
}
