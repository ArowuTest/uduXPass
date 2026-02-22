import React from 'react';
import { Calendar, Users, Settings, Plus } from 'lucide-react';

const IconTest: React.FC = () => {
  return (
    <div style={{ padding: '20px', background: 'white' }}>
      <h3>Icon Test</h3>
      <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
        <Calendar size={24} color="blue" />
        <Users size={24} color="green" />
        <Settings size={24} color="red" />
        <Plus size={24} color="purple" />
      </div>
      <p>If you see colored icons above, Lucide is working correctly.</p>
    </div>
  );
};

export default IconTest;
