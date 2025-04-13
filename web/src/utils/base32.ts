type Variant = '3548' | '4648' | 'hex' | 'clockwork' | 'maki' | 'wah' | 'crockford' | '';
type CrockFordEncoderInput = bigint | number;
type MultiEncoderInput = Uint8Array | string;
interface ModeArray {
    dic: string;
    re: RegExp;
    name: Variant;
    padding: boolean;
    array?: boolean;
    raw?: boolean;
    split?: number; // CrockFordDecoder
    checksum?: boolean; // CrockFordDecoder
}
interface Base32Options {
    variant?: Variant;
    padding?: boolean;
    array?: boolean;
    raw?: boolean;
    split?: number;
    checksum?: boolean;
}
interface ErrorArray {
    isError: boolean;
    message: string;
}
interface ReturnArray {
    data: string | ArrayBuffer;
    error?: ErrorArray;
}

export class Base32 {
    private _mode: ModeArray = { dic: '', padding: true, re: / /, name: '' };
    private _lastError: ErrorArray = { isError: false, message: '' };

    constructor(options: Base32Options = {}) {
        let mode = this._mode = this.setMode(options.variant);
        if (options.padding !== undefined) {
            if (options.padding === true) {
                mode.padding = true;
            } else {
                mode.padding = false;
            }
        }
        if (options.array !== undefined && options.array) {
            mode.array = true;
        }
        if (options.raw !== undefined && options.raw) {
            mode.raw = true;
        }
        if (mode.name === "crockford") {
            if (options.split !== undefined && options.split) {
                mode.split = parseInt('0' + options.split);
            }
            if (options.checksum !== undefined && options.checksum) {
                mode.checksum = true;
            }
        }
    }

    private setMode(variant: Variant = '4648'): ModeArray {
        switch (variant) {
            case 'hex': // RFC4648_HEX
                return {
                    dic: '0123456789ABCDEFGHIJKLMNOPQRSTUV',
                    padding: true,
                    re: /^[A-V0-9]+$/,
                    name: 'hex'
                }
            case 'maki':
            case 'wah':
            case 'clockwork': // Clockwork Base32
                return {
                    dic: '0123456789ABCDEFGHJKMNPQRSTVWXYZ',
                    padding: false,
                    re: /^[A-TV-Z0-9]+$/,
                    name: 'clockwork'
                }
            case 'crockford':
                return {
                    dic: '0123456789ABCDEFGHJKMNPQRSTVWXYZ',
                    padding: false,
                    re: /^[A-TV-Z0-9]+[*~$=U]?$/,
                    name: 'crockford'
                }
            default: // RFC3548 or RFC4648
        }
        return {
            dic: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567',
            padding: true,
            re: /^[A-Z2-7]+$/,
            name: '4648'
        }
    }

    private crockfordEncoder(input: CrockFordEncoderInput): string {
        let input32 = '';
        let output = '';
        const dic = this._mode.dic;

        if (typeof input === "number") {
            input = Math.floor(input);
        }
        if (typeof input === "number" || typeof input === "bigint") {
            if (input > -1) {
                input32 = input.toString(32);
            }
        }
        if (input32.length < 1) {
            this.setError('Invalid data: input number.');
            console.log("Invalid data: input number.");
            return '';
        }

        (input32.split('')).map(index => {
            output += dic[parseInt(index, 32)];
        });

        if (this._mode.checksum) { // check_symbol.length + 1 = 37           
            output += (dic + '*~$=U')[Number(BigInt(input) % BigInt(37))];
        }

        if (this._mode.split && this._mode.split > 0) {
            if (output.length > 0) {
                const reg = new RegExp('(.{1,' + this._mode.split + '})', 'g');
                output = output.match(reg)!.join('-');
            }
        }
        return output;
    }

    private multiEncoder(input: MultiEncoderInput): string {
        if (typeof input !== "object") {
            input = new TextEncoder().encode(input);
        }
        input = new Uint8Array(input);

        const dic = this._mode.dic;
        let output = '';
        let value = 0;
        let offset = 0;

        for (let i = 0; i < input.byteLength; i++) {
            value = (value << 8) | input[i];
            offset += 8;
            while (offset >= 5) {
                output += dic[(value >>> (offset - 5)) & 31];
                offset -= 5;
            }
        }

        if (offset > 0) {
            output += dic[(value << (5 - offset)) & 31];
        }

        if (this._mode.padding && output.length % 8) {
            output += '='.repeat(8 - (output.length % 8));
        }
        return output;
    }

    private crockfordDecoder(input: string = '0'): string | ArrayBuffer {
        input = input.toUpperCase().replace(/[-\s]/g, '').replace(/O/g, '0').replace(/[IL]/g, '1');
        if (this._mode.re.test(input) === false) {
            this.setError('Invalid data: input strings.');
            console.log("Invalid data: input strings.");
            input = '0';
        }

        const dic = this._mode.dic;
        const check_symbol = input.slice(-1);
        if (this._mode.checksum) {
            input = input.slice(0, -1);
        }
        const length = input.length;
        const output = new Uint8Array(Math.ceil(length * 5 / 8));
        let outputHexStr = '';
        let index = output.byteLength
        let value = 0;
        let offset = 0;
        const calcValue = () => {
            outputHexStr = (value & 255).toString(16).padStart(2, '0') + outputHexStr;
            output[--index] = value & 255;
            offset -= 8;
            value = value >>> 8;
        }

        for (let i = length - 1; i >= 0; i--) {
            value = value | (dic.indexOf(input[i]) << offset);
            offset += 5;
            if (offset >= 8) {
                calcValue();
            }
        }
        if (value > 0 || input === '0') {
            calcValue();
        }

        if (output.length > 0 && this._mode.checksum) {
            const verify_symbol = (hexStr: string) => { // '01 ... =U'.length = 37
                return (BigInt('0x' + hexStr) % BigInt(37) !== BigInt((dic + '*~$=U').indexOf(check_symbol)));
            };
            if (verify_symbol(outputHexStr)) {
                this.setError('Invalid data: Checksum error.');
                console.log("Invalid data: Checksum error.");
            }
        }

        if (this._lastError.isError) {
            if (this._mode.raw) {
                return new Uint8Array(1);
            }
            outputHexStr = '0';
        }
        if (this._mode.raw) {
            return output;
        }
        return '0x' + (outputHexStr.replace(/(^0+)(?!$)/, ''));
    }

    private multiDecoder(input: string = ''): string | ArrayBuffer {
        input = input.toUpperCase().replace(/\s|=/g, '');
        if (this._mode.name === 'clockwork') {
            input = input.replace(/O/g, '0').replace(/[IL]/g, '1');
        }
        if (this._mode.re.test(input) === false) {
            this.setError('Invalid data: Input strings.');
            console.log("Invalid data: Input strings.");
            input = '';
        }

        const dic = this._mode.dic;
        const length = input.length;
        const output = new Uint8Array(length * 5 / 8);

        let value = 0;
        let index = 0;
        let offset = 0;

        for (let i = 0; i < length; i++) {
            value = (value << 5) | dic.indexOf(input[i]);
            offset += 5;

            if (offset >= 8) {
                output[index++] = (value >>> (offset - 8)) & 255;
                offset -= 8;
            }
        }

        if (this._mode.raw) {
            return output;
        }
        return new TextDecoder().decode(output.buffer);
    }

    private returnArray(data: string | ArrayBuffer): ReturnArray {
        let ret: ReturnArray = { data: data };
        if (this._lastError.isError) {
            ret.error = this._lastError;
        }
        return ret;
    }

    private setError(message: string): void {
        this._lastError = { isError: !0, message: message };
    }

    private resetError(): void {
        this._lastError = { isError: !1, message: '' };
    }

    public encode(input: MultiEncoderInput | CrockFordEncoderInput): string | ReturnArray {
        this.resetError();
        let data: string;
        if (this._mode.name === 'crockford') {
            data = this.crockfordEncoder(<CrockFordEncoderInput>input);
        } else {
            data = this.multiEncoder(<MultiEncoderInput>input);
        }
        if (this._mode.array) {
            return this.returnArray(data);
        }
        return data;
    };

    public decode(input: string): string | ArrayBuffer | ReturnArray {
        this.resetError();
        let data: string | ArrayBuffer;
        if (this._mode.name === 'crockford') {
            data = this.crockfordDecoder(input);
        } else {
            data = this.multiDecoder(input);
        }
        if (this._mode.array) {
            return this.returnArray(data);
        }
        return data;
    };

    public lasterror(): ErrorArray {
        return this._lastError;
    }

}