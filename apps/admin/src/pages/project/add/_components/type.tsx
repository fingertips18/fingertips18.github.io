import type { Control, FieldErrors, FieldValues, Path } from 'react-hook-form';

import { Dropdown } from '@/components/common/dropdown';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { SelectItem } from '@/components/shadcn/select';
import { ProjectType } from '@/types/project';

interface TypeProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  errors: FieldErrors<T>;
  disabled?: boolean;
}

export function Type<T extends FieldValues>({
  control,
  name,
  errors,
  disabled,
}: TypeProps<T>) {
  const hasError = !!errors[name];

  return (
    <FormField
      control={control}
      name={name}
      disabled={disabled}
      render={({ field }) => {
        const { onChange, ...fields } = field;

        return (
          <FormItem className='w-full'>
            <FormLabel className='w-fit'>Type</FormLabel>
            <FormDescription>
              Select the main platform or category your project belongs to.
            </FormDescription>
            <FormControl>
              <Dropdown
                label='Project Type'
                placeholder='Select a type'
                onValueChange={onChange}
                {...fields}
                hasError={hasError}
              >
                {Object.values(ProjectType).map((t) => (
                  <SelectItem key={t} value={t} className='capitalize'>
                    {t}
                  </SelectItem>
                ))}
              </Dropdown>
            </FormControl>
            <FormMessage />
          </FormItem>
        );
      }}
    />
  );
}
