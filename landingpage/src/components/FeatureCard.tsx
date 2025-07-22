
import React from 'react';

interface FeatureCardProps {
  title: string;
  description: string;
  icon: React.ReactNode;
  className?: string;
}

const FeatureCard: React.FC<FeatureCardProps> = ({ title, description, icon, className = '' }) => {
  return (
    <div className={`p-6 rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm hover:border-terminal-highlight transition-all duration-300 ${className}`}>
      <div className="h-10 w-10 rounded-full bg-terminal-highlight/10 flex items-center justify-center mb-4 text-terminal-highlight">
        {icon}
      </div>
      <h3 className="text-lg font-bold mb-2 text-white">{title}</h3>
      <p className="text-gray-400">{description}</p>
    </div>
  );
};

export default FeatureCard;
