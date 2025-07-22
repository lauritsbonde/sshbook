
import React from 'react';
import { Button } from "@/components/ui/button";
import { Github } from 'lucide-react';

const Header: React.FC = () => {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-black/80 backdrop-blur-sm border-b border-gray-800">
      <div className="container mx-auto flex items-center justify-between py-4">
        <div className="flex items-center">
          <span className="text-terminal-highlight font-mono font-bold text-xl">✨ sshbook</span>
        </div>
        
        <nav className="hidden md:flex items-center space-x-6">
          <a href="#features" className="text-gray-300 hover:text-white transition">Features</a>
          <a href="#why" className="text-gray-300 hover:text-white transition">Why sshbook?</a>
          <a href="https://github.com" target="_blank" rel="noopener noreferrer" className="text-gray-300 hover:text-white transition flex items-center gap-1">
            <Github size={18} />
            <span>GitHub</span>
          </a>
        </nav>
        
        <div>
          <Button className="bg-terminal-highlight text-gray-900 hover:bg-terminal-highlight/90">Download</Button>
        </div>
      </div>
    </header>
  );
};

export default Header;
