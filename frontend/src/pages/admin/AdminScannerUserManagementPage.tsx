import React, { useState, useEffect } from 'react';
import { scannerUsersAPI, eventsAPI } from '@/services/api';
import { toast } from "@/components/ui/toaster";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Switch } from '@/components/ui/switch';
import { Loader2, UserPlus, RefreshCw, Scan, Shield, Settings, User, CheckCircle, XCircle } from 'lucide-react';
import { AdminUser, Event } from '@/types/api';

const AdminScannerUserManagementPage: React.FC = () => {
  
  const [scannerUsers, setScannerUsers] = useState<AdminUser[]>([]);
  const [events, setEvents] = useState<Event[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [showCreateUserDialog, setShowCreateUserDialog] = useState<boolean>(false);
  const [newUserData, setNewUserData] = useState({
    username: '',
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    role: 'scanner',
    is_active: true,
    assigned_events: []
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [usersResponse, eventsResponse] = await Promise.all([
        scannerUsersAPI.getAll(),
        eventsAPI.getAll()
      ]);

      if (usersResponse.success) {
        setScannerUsers(usersResponse.data.data || []);
      } else {
        toast({ title: 'Error', description: usersResponse.error || 'Failed to fetch scanner users.', variant: 'destructive' });
      }

      if (eventsResponse.success) {
        setEvents(eventsResponse.data || []);
      } else {
        toast({ title: 'Error', description: eventsResponse.error || 'Failed to fetch events.', variant: 'destructive' });
      }
    } catch (error) {
      toast({ title: 'Error', description: 'An unexpected error occurred.', variant: 'destructive' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      const response = await scannerUsersAPI.create(newUserData);
      if (response.success) {
        toast({ title: 'Success', description: 'Scanner user created successfully.' });
        setShowCreateUserDialog(false);
        fetchData();
      } else {
        toast({ title: 'Error', description: response.error || 'Failed to create scanner user.', variant: 'destructive' });
      }
    } catch (error) {
      toast({ title: 'Error', description: 'An unexpected error occurred.', variant: 'destructive' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const getRoleBadge = (role: string): JSX.Element => {
    const roleConfig: Record<string, { color: string; label: string; icon: React.ComponentType<any> }> = {
      scanner: { color: 'bg-blue-100 text-blue-800', label: 'Scanner', icon: Scan },
      supervisor: { color: '" style={{ background: "rgba(245,158,11,0.1)" }} text-purple-800', label: 'Supervisor', icon: Shield },
      admin: { color: 'bg-red-100 text-red-800', label: 'Admin', icon: Settings }
    };
    const config = roleConfig[role] || { color: 'bg-opacity-0" style={{ background: "var(--brand-surface)" }} font-semibold', label: role, icon: User };
    const Icon = config.icon;
    return (
      <Badge className={`${config.color} flex items-center gap-1`}>
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    );
  };

  const getStatusBadge = (isActive: boolean): JSX.Element => {
    return (
      <Badge className={isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}>
        {isActive ? <CheckCircle className="h-3 w-3 mr-1" /> : <XCircle className="h-3 w-3 mr-1" />}
        {isActive ? 'Active' : 'Inactive'}
      </Badge>
    );
  };

  if (isLoading) {
    return <div className="flex justify-center items-center h-64"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Scanner User Management</h1>
          <p className="text-muted-foreground">Manage scanner user accounts, passwords, and event assignments</p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" onClick={fetchData}><RefreshCw className="h-4 w-4 mr-2" />Refresh</Button>
          <Dialog open={showCreateUserDialog} onOpenChange={setShowCreateUserDialog}>
            <DialogTrigger asChild>
              <Button><UserPlus className="h-4 w-4 mr-2" />Add Scanner User</Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Create New Scanner User</DialogTitle>
                <DialogDescription>Add a new scanner user account with role and event assignments</DialogDescription>
              </DialogHeader>
              <form onSubmit={handleCreateUser} className="space-y-4">
                {/* Form fields... */}
                <Button type="submit" disabled={isSubmitting}>{isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Create User'}</Button>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Scanner Users</CardTitle>
          <CardDescription>A list of all scanner users in the system.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">Username</th>
                  <th className="text-left p-4">Email</th>
                  <th className="text-left p-4">Role</th>
                  <th className="text-left p-4">Status</th>
                  <th className="text-left p-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {scannerUsers.map((user) => (
                  <tr key={user.id}>
                    <td className="p-4">{user.username}</td>
                    <td className="p-4">{user.email}</td>
                    <td className="p-4">{getRoleBadge(user.role)}</td>
                    <td className="p-4">{getStatusBadge(user.is_active)}</td>
                    <td className="p-4">{/* Action buttons... */}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default AdminScannerUserManagementPage;

