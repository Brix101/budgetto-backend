import { queryClient } from "@/lib/react-query-client";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Toaster } from "sonner";

export const AppProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        {children}
        <ReactQueryDevtools initialIsOpen={false} />
        <Toaster richColors closeButton />
      </QueryClientProvider>
    </>
  );
};
