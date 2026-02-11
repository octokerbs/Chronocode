import { MarketingPanel } from "@/components/login/marketing-panel";
import { LoginForm } from "@/components/login/login-form";

export default function LoginPage() {
  return (
    <div className="flex h-screen">
      {/* Left: Marketing panel (60%) */}
      <div className="hidden w-3/5 border-r bg-muted/30 lg:block">
        <MarketingPanel />
      </div>

      {/* Right: Login form (40%) */}
      <div className="flex w-full items-center justify-center lg:w-2/5">
        <LoginForm />
      </div>
    </div>
  );
}
