
import React from 'react';
import { Button } from "@/components/ui/button";
import Header from '@/components/Header';
import Footer from '@/components/Footer';
import TerminalDemo from '@/components/TerminalDemo';
import FeatureCard from '@/components/FeatureCard';
import { Archive, Zap, BookOpen, Smartphone, Github } from 'lucide-react';

const Index = () => {
  return (
    <div className="min-h-screen bg-black text-white flex flex-col">
      <Header />
      
      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4">
        <div className="container mx-auto grid md:grid-cols-2 gap-12 items-center">
          <div className="animate-fade-in">
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6">
              Your SSH connections,<br />
              <span className="text-terminal-highlight">organized.</span> Everywhere.
            </h1>
            <p className="text-lg md:text-xl text-gray-400 mb-8 max-w-lg">
              sshbook is a lightweight, developer-focused CLI tool that helps you organize, label, and connect to all your servers — right from your terminal.
            </p>
            <div className="flex flex-col sm:flex-row gap-4">
              <Button className="bg-terminal-highlight hover:bg-terminal-highlight/90 text-gray-900 px-6 py-6 text-lg">
                Get Started
              </Button>
              <Button 
                variant="outline" 
                className="border-terminal-highlight/50 hover:bg-gray-800 px-6 py-6 text-lg text-terminal-highlight flex items-center gap-2"
              >
                <Github size={20} />
                View on GitHub
              </Button>
            </div>
          </div>
          
          <div className="animate-fade-in">
            <TerminalDemo className="w-full max-w-md mx-auto" />
          </div>
        </div>
      </section>
      
      {/* Features Section */}
      <section id="features" className="py-20 px-4 bg-gradient-to-b from-black to-gray-950">
        <div className="container mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">Simplify your SSH workflow</h2>
            <p className="text-gray-400 max-w-2xl mx-auto">
              No more digging through messy .ssh/config files or struggling to remember hostnames.
              With sshbook, your servers are grouped, searchable, and always just a few keystrokes away.
            </p>
          </div>
          
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            <FeatureCard 
              title="Group and label" 
              description="Organize your SSH connections with custom groups and labels for easy access."
              icon={<Archive size={20} />}
            />
            <FeatureCard 
              title="Instant connect" 
              description="Connect to your servers with just a few keystrokes right from your terminal."
              icon={<Zap size={20} />}
            />
            <FeatureCard 
              title="Simple config" 
              description="Simple, clean configuration — not a bloated manager."
              icon={<BookOpen size={20} />}
            />
            <FeatureCard 
              title="Sync across devices" 
              description="(Pro) Access your connections from anywhere, including your phone."
              icon={<Smartphone size={20} />}
            />
          </div>
        </div>
      </section>
      
      {/* Why sshbook Section */}
      <section id="why" className="py-20 px-4 bg-gray-950">
        <div className="container mx-auto">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-3xl md:text-4xl font-bold mb-12 flex items-center justify-center">
              <span className="bg-terminal-highlight/20 text-terminal-highlight px-4 py-1 rounded-lg mr-3">🛠</span> 
              Why sshbook?
            </h2>
            <blockquote className="text-2xl md:text-3xl font-light italic mb-8 text-gray-300">
              "Think of it as your personal address book for servers — simple, fast, and made for people who love the command line."
            </blockquote>
            <div className="h-px w-16 bg-terminal-highlight mx-auto mb-8"></div>
            <p className="text-gray-400">
              Built for developers, sysadmins, and anyone who lives in their terminal.
            </p>
          </div>
        </div>
      </section>
      
      {/* CTA Section */}
      <section className="py-20 px-4 bg-terminal-highlight/10 border-y border-terminal-highlight/20">
        <div className="container mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold mb-6">
            Keep your servers a <span className="text-terminal-highlight">keystroke away</span>
          </h2>
          <p className="text-gray-400 max-w-2xl mx-auto mb-8">
            Very fast, zero-bloat CLI. Local JSON storage. Simple but expandable.
          </p>
          <Button className="bg-terminal-highlight hover:bg-terminal-highlight/90 text-gray-900 px-8 py-6 text-lg">
            Download sshbook
          </Button>
        </div>
      </section>
      
      {/* Terminal Command examples */}
      <section className="py-20 px-4 bg-gray-950">
        <div className="container mx-auto">
          <div className="max-w-3xl mx-auto bg-terminal-DEFAULT rounded-lg p-6 font-mono text-terminal-text">
            <div className="mb-4">
              <span className="terminal-prompt">$</span> sshbook add web-server
              <div className="pl-4">
                <div><span className="opacity-75">&gt; Host:</span> web-01.example.com</div>
                <div><span className="opacity-75">&gt; User:</span> admin</div>
                <div><span className="opacity-75">&gt; Group:</span> production</div>
                <div><span className="text-terminal-highlight">&gt; SSH connection added successfully!</span></div>
              </div>
            </div>
            
            <div className="mb-4">
              <span className="terminal-prompt">$</span> sshbook list
              <div className="pl-4 pt-1">
                <div><span className="opacity-75">📖 sshbook connections:</span></div>
                <div className="pl-2"><span className="opacity-75">├─ 📁 development</span></div>
                <div className="pl-2"><span className="opacity-75">│  └─ 🖥️  dev-server</span></div>
                <div className="pl-2"><span className="opacity-75">└─ 📁 production</span></div>
                <div className="pl-2"><span className="opacity-75">   └─ 🖥️  web-server</span></div>
              </div>
            </div>
            
            <div>
              <span className="terminal-prompt">$</span> sshbook connect web-server
              <div className="pl-4">
                <div><span className="text-terminal-highlight">✨ Connecting to web-server...</span></div>
              </div>
            </div>
          </div>
        </div>
      </section>
      
      <Footer />
    </div>
  );
};

export default Index;
