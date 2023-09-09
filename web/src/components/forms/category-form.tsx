import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import {
  Category,
  createCategorySchema,
  updateCategorySchema,
} from "@/lib/validations/category";
import {
  createCategory,
  deleteCategory,
  updateCategory,
} from "@/services/category.service";
import { useAuth0 } from "@auth0/auth0-react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { Icons } from "@/components/icons";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import { useBoundStore } from "@/lib/store";
import { useEffect } from "react";

type CreateInputs = z.infer<typeof createCategorySchema>;
type UpdateInputs = z.infer<typeof updateCategorySchema>;

export function CategoryCreateDialog() {
  const auth = useAuth0();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { mode, setMode } = useBoundStore((state) => state.category);

  const form = useForm<CreateInputs>({
    resolver: zodResolver(createCategorySchema),
    defaultValues: {
      name: "",
      note: "",
    },
  });

  const { mutate, isLoading } = useMutation({
    mutationFn: createCategory,
    onSuccess: (response) => {
      queryClient.setQueriesData([QUERY_CATEGORIES_KEY], (prev: unknown) => {
        const categories = prev as Category[];
        return [response, ...categories];
      });

      handleCancelClick();
      toast({
        title: "Created successfully",
        description: `category ${response.name} created successfully`,
      });
    },
    onError: (error) => {
      console.log({ error });
    },
  });

  function onSubmit(data: CreateInputs) {
    mutate({ auth, category: data });
  }

  function handleCancelClick() {
    setMode({ mode: "view" });
    form.reset();
  }

  return (
    <AlertDialog open={mode === "create"}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create category</AlertDialogTitle>
          <AlertDialogDescription>
            Create your new custom category
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Separator />
        <div>
          <Form {...form}>
            <form
              className="grid gap-4"
              onSubmit={(...args) => void form.handleSubmit(onSubmit)(...args)}
            >
              <FormField
                control={form.control}
                name="name"
                disabled={isLoading}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="category name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="note"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="category description"
                        className="resize-none"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter>
                <Button
                  type="button"
                  disabled={isLoading}
                  variant={"outline"}
                  onClick={handleCancelClick}
                >
                  Cancel
                </Button>
                <Button disabled={isLoading}>
                  {isLoading && (
                    <Icons.spinner
                      className="mr-2 h-4 w-4 animate-spin"
                      aria-hidden="true"
                    />
                  )}
                  Create
                  <span className="sr-only">Create</span>
                </Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export function CategoryDeleteDialog() {
  const auth = useAuth0();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { mode, category, setMode } = useBoundStore((state) => state.category);

  const { mutate, isLoading } = useMutation({
    mutationFn: deleteCategory,
    onSuccess: (response) => {
      queryClient.setQueriesData([QUERY_CATEGORIES_KEY], (prev: unknown) => {
        const categories = prev as Category[];
        return categories.filter((item) => item.id !== category?.id);
      });

      handleCancelClick();
      toast({
        title: "Deleted successfully",
        description: response.message,
      });
    },
    onError: (error) => {
      console.log({ error });
    },
  });

  function handleCancelClick() {
    setMode({ mode: "view" });
  }

  function handleDeleteClick() {
    mutate({ auth, id: category?.id ?? 0 });
  }

  return (
    <AlertDialog open={mode === "delete"}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete your
            cateory and remove your data from our servers.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <Button
            variant={"destructive"}
            disabled={isLoading}
            onClick={handleDeleteClick}
          >
            {isLoading && (
              <Icons.spinner
                className="mr-2 h-4 w-4 animate-spin"
                aria-hidden="true"
              />
            )}
            Continue
          </Button>

          <AlertDialogCancel disabled={isLoading} onClick={handleCancelClick}>
            Cancel
          </AlertDialogCancel>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export function CategoryUpdateDialog() {
  const auth = useAuth0();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { mode, category, setMode } = useBoundStore((state) => state.category);

  const form = useForm<UpdateInputs>({
    resolver: zodResolver(updateCategorySchema),
    defaultValues: {
      id: category?.id,
      name: category?.name,
      note: category?.note ?? "",
    },
  });

  const { mutate, isLoading } = useMutation({
    mutationFn: updateCategory,
    onSuccess: (response) => {
      queryClient.setQueriesData([QUERY_CATEGORIES_KEY], (prev: unknown) => {
        const categories = prev as Category[];
        return categories.map((item) => {
          if (item.id === response.id) {
            return response;
          }
          return item;
        });
      });

      handleCancelClick();
      toast({
        title: "Updated successfully",
        description: `category ${response.name} updated successfully`,
      });
    },
    onError: (error) => {
      console.log({ error });
    },
  });

  function onSubmit(data: UpdateInputs) {
    mutate({ auth, category: data });
  }

  function handleCancelClick() {
    setMode({ mode: "view" });
    form.reset();
  }

  useEffect(() => {
    if (mode === "update") {
      form.reset({
        id: category?.id,
        name: category?.name,
        note: category?.note ?? "",
      });
    }
  }, [mode, category]);

  return (
    <AlertDialog open={mode === "update"}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Update category</AlertDialogTitle>
          <AlertDialogDescription>Update your category</AlertDialogDescription>
        </AlertDialogHeader>
        <Separator />
        <div>
          <Form {...form}>
            <form
              className="grid gap-4"
              onSubmit={(...args) => void form.handleSubmit(onSubmit)(...args)}
            >
              <FormField
                control={form.control}
                name="name"
                disabled={isLoading}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="category name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="note"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="category description"
                        className="resize-none"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter>
                <Button
                  type="button"
                  disabled={isLoading}
                  variant={"outline"}
                  onClick={handleCancelClick}
                >
                  Cancel
                </Button>
                <Button disabled={isLoading}>
                  {isLoading && (
                    <Icons.spinner
                      className="mr-2 h-4 w-4 animate-spin"
                      aria-hidden="true"
                    />
                  )}
                  Update
                  <span className="sr-only">Update</span>
                </Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  );
}
