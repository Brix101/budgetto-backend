import {
  ArrowDown,
  ArrowRightLeft,
  ArrowUp,
  BadgeDollarSign,
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  CreditCard,
  DollarSign,
  Eye,
  EyeOff,
  FileTerminal,
  LogOut,
  Moon,
  MoreHorizontal,
  MoreVertical,
  MoveLeft,
  Search,
  Settings,
  ShoppingBag,
  Sun,
  User,
  X,
} from "lucide-react";

type IconProps = React.HTMLAttributes<SVGElement>;

export const Icons = {
  arrowUp: ArrowUp,
  arrowDown: ArrowDown,
  close: X,
  view: Eye,
  hide: EyeOff,
  moveLeft: MoveLeft,
  user: User,
  terminal: FileTerminal,
  settings: Settings,
  logout: LogOut,
  store: ShoppingBag,
  billing: CreditCard,
  dollarSign: DollarSign,
  transaction: ArrowRightLeft,
  account: BadgeDollarSign,
  chevronsLeft: ChevronsLeft,
  chevronLeft: ChevronLeft,
  chevronsRight: ChevronsRight,
  chevronRight: ChevronRight,
  moon: Moon,
  sun: Sun,
  verticalThreeDots: MoreVertical,
  horizontalThreeDots: MoreHorizontal,
  search: Search,
  category: (props: IconProps) => (
    <svg
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M10 13H3C2.73478 13 2.48043 13.1054 2.29289 13.2929C2.10536 13.4804 2 13.7348 2 14V21C2 21.2652 2.10536 21.5196 2.29289 21.7071C2.48043 21.8946 2.73478 22 3 22H10C10.2652 22 10.5196 21.8946 10.7071 21.7071C10.8946 21.5196 11 21.2652 11 21V14C11 13.7348 10.8946 13.4804 10.7071 13.2929C10.5196 13.1054 10.2652 13 10 13ZM9 20H4V15H9V20ZM21 2H14C13.7348 2 13.4804 2.10536 13.2929 2.29289C13.1054 2.48043 13 2.73478 13 3V10C13 10.2652 13.1054 10.5196 13.2929 10.7071C13.4804 10.8946 13.7348 11 14 11H21C21.2652 11 21.5196 10.8946 21.7071 10.7071C21.8946 10.5196 22 10.2652 22 10V3C22 2.73478 21.8946 2.48043 21.7071 2.29289C21.5196 2.10536 21.2652 2 21 2ZM20 9H15V4H20V9ZM10 2H3C2.73478 2 2.48043 2.10536 2.29289 2.29289C2.10536 2.48043 2 2.73478 2 3V10C2 10.2652 2.10536 10.5196 2.29289 10.7071C2.48043 10.8946 2.73478 11 3 11H10C10.2652 11 10.5196 10.8946 10.7071 10.7071C10.8946 10.5196 11 10.2652 11 10V3C11 2.73478 10.8946 2.48043 10.7071 2.29289C10.5196 2.10536 10.2652 2 10 2ZM9 9H4V4H9V9Z"
        fill="currentColor"
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M21.9777 17.947C21.8743 18.9829 21.415 19.9508 20.678 20.686C20.2593 21.1036 19.7624 21.4346 19.2157 21.6601C18.6691 21.8855 18.0833 22.001 17.492 22C16.451 21.9981 15.4428 21.6354 14.6393 20.9736C13.8357 20.3118 13.2864 19.3919 13.085 18.3706C12.8836 17.3493 13.0426 16.2897 13.5347 15.3724C14.0269 14.4551 14.8219 13.7368 15.7843 13.3399C16.7466 12.943 17.8168 12.892 18.8126 13.1957C19.8083 13.4993 20.6679 14.1388 21.2451 15.0051C21.8222 15.8715 22.0811 16.9112 21.9777 17.947ZM17.5 20C18.8807 20 20 18.8807 20 17.5C20 16.1193 18.8807 15 17.5 15C16.1193 15 15 16.1193 15 17.5C15 18.8807 16.1193 20 17.5 20Z"
        fill="currentColor"
      />
    </svg>
  ),
  logo: (props: IconProps) => (
    <svg
      width="60"
      height="60"
      viewBox="0 0 60 60"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M45.7993 20.0996C45.2582 14.4749 40.4023 12.59 34.2659 12.0524V4.24707H29.5175V11.844C28.2707 11.844 26.9937 11.8676 25.7268 11.8944V4.24707H20.9784L20.975 12.0456C19.9467 12.0658 18.9351 12.086 17.9505 12.086V12.0624L11.4008 12.0591V17.1327C11.4008 17.1327 14.9092 17.0655 14.8487 17.1293C16.7743 17.1293 17.3994 18.2448 17.5809 19.2092V28.0997C17.7153 28.0997 17.8867 28.1064 18.0816 28.1333H17.5809L17.5775 40.5888C17.4935 41.1936 17.1373 42.158 15.793 42.1613C15.8535 42.2151 12.3418 42.1613 12.3418 42.1613L11.3975 47.833H17.5809C18.7302 47.833 19.8627 47.8532 20.9716 47.8599L20.975 55.7525H25.7201V47.9439C27.0206 47.9708 28.2808 47.9808 29.5141 47.9808L29.5108 55.7525H34.2592V47.8767C42.2438 47.4197 47.8391 45.4071 48.5314 37.9076C49.0926 31.8696 46.253 29.1716 41.7196 28.0829C44.4786 26.6852 46.2026 24.2156 45.7993 20.0996V20.0996ZM39.1521 36.9735C39.1521 42.8703 29.0537 42.2016 25.831 42.2016V31.742C29.0537 31.7487 39.1521 30.8247 39.1521 36.9735V36.9735ZM36.9409 22.2197C36.9409 27.5856 28.5127 26.9573 25.831 26.9607V17.4788C28.516 17.4788 36.9443 16.6253 36.9409 22.2197Z"
        fill="currentColor"
      />
      <path
        d="M17.5869 27.8916H18.2926V28.4124H17.5869V27.8916Z"
        fill="currentColor"
      />
    </svg>
  ),
  spinner: (props: IconProps) => (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      {...props}
    >
      <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
  ),
};
