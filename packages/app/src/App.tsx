import { PlatformLayout } from "@/layouts/Platform/Platform";
import { PackCalculator } from "@/pages/PackCalculator/PackCalculator";

export default function App() {
  return (
    <PlatformLayout title="Gymshark Pack Calculator">
      <PackCalculator />
    </PlatformLayout>
  );
}
