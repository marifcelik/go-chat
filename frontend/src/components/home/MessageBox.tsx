import { useEffect, useState } from 'react'
import { SendIcon } from 'lucide-react'
import { Input } from '../ui/input'
import { Button } from '../ui/button'

export default function MessageBox({
	onSubmit,
	value,
	setValue
}: {
	onSubmit: (e: React.FormEvent<HTMLFormElement>) => void
	value: string
	setValue: (value: React.SetStateAction<string>) => void
}) {
	const [disabled, setDisabled] = useState(true)

	useEffect(() => {
		setDisabled(value.length === 0)
	}, [value])

	return (
		<div id="messageBox" className="border-t h-16 w-full">
			<form onSubmit={onSubmit} className="flex w-full items-center space-x-2 p-3">
				<Input
					name="message"
					placeholder="Type your message..."
					className="flex-1"
					autoComplete="off"
					value={value}
					onChange={e => setValue(e.target.value)}
				/>
				<Button type="submit" size="icon" disabled={disabled}>
					<span className="sr-only">Send</span>
					<SendIcon className="h-4 w-4" />
				</Button>
			</form>
		</div>
	)
}