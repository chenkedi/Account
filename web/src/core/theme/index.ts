import { extendTheme, type ThemeConfig } from '@chakra-ui/react';
import { colors } from './colors';

const config: ThemeConfig = {
  initialColorMode: 'light',
  useSystemColorMode: false,
};

export const appTheme = extendTheme({
  config,
  colors: {
    brand: colors.brand,
    income: colors.income,
    expense: colors.expense,
    transfer: colors.transfer,
    gray: colors.gray,
  },
  fonts: {
    heading:
      '-apple-system, BlinkMacSystemFont, "SF Pro Display", "Segoe UI", Roboto, sans-serif',
    body: '-apple-system, BlinkMacSystemFont, "SF Pro Text", "Segoe UI", Roboto, sans-serif',
  },
  fontSizes: {
    xs: '0.75rem',
    sm: '0.875rem',
    md: '1rem',
    lg: '1.125rem',
    xl: '1.25rem',
    '2xl': '1.5rem',
    '3xl': '1.875rem',
    '4xl': '2.25rem',
  },
  radii: {
    none: '0',
    sm: '0.375rem',
    md: '0.5rem',
    lg: '0.75rem',
    xl: '1rem',
    '2xl': '1.5rem',
    '3xl': '2rem',
    full: '9999px',
  },
  shadows: {
    xs: '0 1px 2px 0px rgba(0, 0, 0, 0.05)',
    sm: '0 1px 3px 0px rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
    md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
    lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
    xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
    '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
    soft: '0 2px 20px rgba(0, 0, 0, 0.08)',
    card: '0 4px 20px rgba(99, 102, 241, 0.08)',
    button: '0 4px 12px rgba(99, 102, 241, 0.3)',
  },
  space: {
    px: '1px',
    0.5: '0.125rem',
    1: '0.25rem',
    1.5: '0.375rem',
    2: '0.5rem',
    2.5: '0.625rem',
    3: '0.75rem',
    3.5: '0.875rem',
    4: '1rem',
    5: '1.25rem',
    6: '1.5rem',
    7: '1.75rem',
    8: '2rem',
    9: '2.25rem',
    10: '2.5rem',
    12: '3rem',
    14: '3.5rem',
    16: '4rem',
    20: '5rem',
    24: '6rem',
    28: '7rem',
    32: '8rem',
    36: '9rem',
    40: '10rem',
    44: '11rem',
    48: '12rem',
    52: '13rem',
    56: '14rem',
    60: '15rem',
    64: '16rem',
    72: '18rem',
    80: '20rem',
    96: '24rem',
  },
  styles: {
    global: {
      'html, body, #root': {
        height: '100%',
        bg: 'gray.50',
        color: 'gray.900',
        fontFeatureSettings: '"cv02", "cv03", "cv04", "cv11"',
      },
      '*::placeholder': {
        color: 'gray.400',
      },
      '*': {
        borderColor: 'gray.200',
      },
    },
  },
  components: {
    Button: {
      defaultProps: {
        colorScheme: 'brand',
        size: 'md',
      },
      baseStyle: {
        borderRadius: 'xl',
        fontWeight: '600',
        _focus: {
          boxShadow: 'none',
        },
      },
      variants: {
        solid: {
          bg: 'brand.500',
          color: 'white',
          boxShadow: 'button',
          _hover: {
            bg: 'brand.600',
            boxShadow: 'md',
            _disabled: {
              bg: 'brand.500',
            },
          },
          _active: {
            bg: 'brand.700',
          },
        },
        outline: {
          borderWidth: '2px',
          borderColor: 'brand.500',
          color: 'brand.500',
          bg: 'transparent',
          _hover: {
            bg: 'brand.50',
          },
          _active: {
            bg: 'brand.100',
          },
        },
        ghost: {
          color: 'gray.600',
          _hover: {
            bg: 'gray.100',
            color: 'gray.900',
          },
          _active: {
            bg: 'gray.200',
          },
        },
        soft: {
          bg: 'brand.100',
          color: 'brand.700',
          _hover: {
            bg: 'brand.200',
          },
          _active: {
            bg: 'brand.300',
          },
        },
      },
      sizes: {
        lg: {
          h: '14',
          px: '6',
          fontSize: 'lg',
        },
        md: {
          h: '12',
          px: '5',
          fontSize: 'md',
        },
        sm: {
          h: '10',
          px: '4',
          fontSize: 'sm',
        },
        xs: {
          h: '8',
          px: '3',
          fontSize: 'xs',
        },
      },
    },
    Input: {
      defaultProps: {
        focusBorderColor: 'brand.500',
      },
      baseStyle: {
        field: {
          fontWeight: '500',
          _placeholder: {
            color: 'gray.400',
          },
        },
      },
      variants: {
        outline: {
          field: {
            borderWidth: '2px',
            borderColor: 'gray.200',
            bg: 'white',
            borderRadius: 'xl',
            _hover: {
              borderColor: 'gray.300',
            },
            _focus: {
              borderColor: 'brand.500',
              boxShadow: '0 0 0 4px rgba(99, 102, 241, 0.1)',
            },
            _invalid: {
              borderColor: 'expense.500',
              _focus: {
                boxShadow: '0 0 0 4px rgba(239, 68, 68, 0.1)',
              },
            },
          },
        },
        filled: {
          field: {
            bg: 'gray.100',
            borderRadius: 'xl',
            _hover: {
              bg: 'gray.200',
            },
            _focus: {
              bg: 'white',
              borderColor: 'brand.500',
              boxShadow: '0 0 0 4px rgba(99, 102, 241, 0.1)',
            },
          },
        },
      },
      sizes: {
        lg: {
          field: {
            h: '14',
            fontSize: 'lg',
            px: '5',
          },
        },
        md: {
          field: {
            h: '12',
            fontSize: 'md',
            px: '4',
          },
        },
        sm: {
          field: {
            h: '10',
            fontSize: 'sm',
            px: '3',
          },
        },
      },
    },
    Card: {
      baseStyle: {
        container: {
          borderRadius: '2xl',
          borderWidth: '1px',
          borderColor: 'gray.100',
          bg: 'white',
          boxShadow: 'soft',
        },
      },
      variants: {
        elevated: {
          container: {
            boxShadow: 'card',
            borderWidth: '0',
          },
        },
        outline: {
          container: {
            bg: 'white',
            boxShadow: 'none',
          },
        },
        glass: {
          container: {
            bg: 'whiteAlpha.800',
            backdropFilter: 'blur(12px)',
            borderColor: 'whiteAlpha.400',
          },
        },
      },
      defaultProps: {
        variant: 'outline',
      },
    },
    Heading: {
      baseStyle: {
        fontWeight: '700',
        letterSpacing: '-0.02em',
      },
    },
    Text: {
      baseStyle: {
        color: 'gray.700',
      },
    },
    Badge: {
      defaultProps: {
        colorScheme: 'brand',
      },
      baseStyle: {
        borderRadius: 'full',
        fontWeight: '600',
        px: '2.5',
        py: '1',
      },
    },
    Tag: {
      baseStyle: {
        container: {
          borderRadius: 'full',
          fontWeight: '500',
        },
      },
    },
  },
});
